package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/steveoc64/go-cmms/shared"
)

type UtilRPC struct{}

// Do a simple database backup
func (u *UtilRPC) Backup(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)
	*result = ""

	if conn.UserRole == "Admin" {
		out, err := exec.Command("../scripts/cmms-backup.sh").Output()
		if err != nil {
			log.Println(err)
			*result = err.Error()
			return nil
		}
		*result = string(out)
	}

	logger(start, "Util.Backup",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		*result,
		channel, conn.UserID, "", 0, false)

	return nil
}

// Run a top command and return the results
func (u *UtilRPC) Top(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)
	*result = ""

	if conn.UserRole == "Admin" {
		out, err := exec.Command("../scripts/top.sh").Output()
		if err != nil {
			log.Println(err)
			*result = err.Error()
			return nil
		}
		*result = string(out)
	}

	logger(start, "Util.Top",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		"top",
		channel, conn.UserID, "", 0, false)

	return nil
}

// Tail the output of the server process log file, and return the result
func (u *UtilRPC) Logs(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)
	*result = ""

	if conn.UserRole == "Admin" {
		out, err := exec.Command("../scripts/logs.sh").Output()
		if err != nil {
			log.Println(err)
			*result = err.Error()
			return nil
		}
		*result = string(out)
	}

	logger(start, "Util.Logs",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		"logs",
		channel, conn.UserID, "", 0, false)

	return nil
}

// Patch up tha machine PartClass - Steve only
func (u *UtilRPC) Machine(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)
	*result = ""

	if conn.UserRole == "Admin" && conn.Username == "steve" {

		// For each machine, set the class the same as the name, and report any errors
		r := "Processing Machines\n"

		machines := []shared.Machine{}
		partClass := shared.PartClass{}

		DB.SQL(`select 
			m.*,s.name as site_name
			from machine m
			left join site s on s.id=m.site_id
			where part_class=0
			order by m.id`).
			QueryStructs(&machines)

		patched := 0
		for _, m := range machines {
			siteName := ""
			if m.SiteName != nil {
				siteName = *m.SiteName
			}
			r += fmt.Sprintf("Machine %d: %s (%s)", m.ID, m.Name, siteName)

			err := DB.SQL(`select * from part_class where name=$1`, m.Name).QueryStruct(&partClass)
			if err != nil {
				r += fmt.Sprintf("\n    !! No Matching Part Class !!\n")
				continue
			}

			if partClass.ID == 0 {
				r += fmt.Sprintf("\n    !! No Matching Part Class !!\n")
			} else {
				r += fmt.Sprintf(" = PartClass %d: %s\n", partClass.ID, partClass.Name)
				DB.SQL(`update machine set part_class=$1 where id=$2`, partClass.ID, m.ID).Exec()
				patched++
			}
		}

		r += fmt.Sprintf("\nPatched %d of %d Machines\n", patched, len(machines))
		*result = r
	}

	logger(start, "Util.Machine",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		*result,
		channel, conn.UserID, "part_class", 0, true)

	return nil
}

// Patch up the Parts PartClass field - Steve only
func (u *UtilRPC) Parts(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)

	if conn.UserRole == "Admin" && conn.Username == "steve" {
		// For each part, get the 1st component that its associated
		// with (under the old scheme), and from there get the machine.
		//
		// The machine then tells us which partclass to use
		r := "Processing Parts\n"

		parts := []shared.Part{}

		DB.SQL(`select * from part
			where class=0
			order by id`).
			QueryStructs(&parts)

		patched := 0
		for _, p := range parts {
			r += fmt.Sprintf("Part %d: %s", p.ID, p.Name)

			// Get the first associated component
			classID := 0
			err := DB.SQL(`select 
				m.part_class
				from component_part x
				left join component c on c.id=x.component_id
				left join machine m on m.id=c.machine_id
				where x.part_id=$1
				limit 1`, p.ID).QueryScalar(&classID)

			if err != nil {
				r += fmt.Sprintf("\n    !! No can find the partclass\n")
				continue
			}

			if classID == 0 {
			} else {
				r += fmt.Sprintf(" = PartClass %d\n", classID)
				DB.SQL(`update part set class=$1 where id=$2`, classID, p.ID).Exec()
				patched++
			}
		}

		r += fmt.Sprintf("\nPatched %d of %d Parts\n", patched, len(parts))
		*result = r
	}

	logger(start, "Util.Parts",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		*result,
		channel, conn.UserID, "part", 0, true)

	return nil
}

// Construct the parts categories for bootstrap
func (u *UtilRPC) Cats(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)

	if conn.UserRole == "Admin" && conn.Username == "steve" {
		// For each part, get the 1st component that its associated
		// with (under the old scheme), and from there get the machine.
		//
		// The machine then tells us which partclass to use
		r := "Processing Cats\n"

		// Truncate the cats
		DB.SQL(`truncate category restart identity`).Exec()

		// Create 1 top level cat for each machine type
		partClasses := []shared.PartClass{}
		DB.SQL(`select * from part_class order by name`).QueryStructs(&partClasses)

		for i, p := range partClasses {
			fmt.Printf("%d: %v\n", i, p)

			cat := shared.Category{
				Name:  p.Name,
				Descr: p.Descr,
			}

			DB.InsertInto("category").
				Whitelist("name", "descr").
				Record(cat).
				Returning("id").
				QueryScalar(&cat.ID)
			fmt.Printf("%v\n", cat)

			// With this new category, stamp ALL parts records with this cat id, where
			// the partclass == selected partclass

			DB.SQL(`update part set category=$1 where class=$2`, cat.ID, p.ID).Exec()

			// Create a sub-category under this one, for each tool in the machine
			machine := shared.Machine{}
			DB.SQL(`select * from machine where part_class=$1 limit 1`, p.ID).
				QueryStruct(&machine)

			components := []shared.Component{}
			DB.SQL(`select * from component where machine_id=$1 order by position`, machine.ID).
				QueryStructs(&components)

			for j, c := range components {
				fmt.Printf("%d: %v\n", j, c)

				subcat := shared.Category{
					ParentID: cat.ID,
					Name:     c.Name,
					Descr:    c.Descr,
				}

				DB.InsertInto("category").
					Whitelist("parent_id", "name", "descr").
					Record(subcat).
					Returning("id").
					QueryScalar(&subcat.ID)
				fmt.Printf("%v\n", subcat)

				// get all parts in this category, and stamp the category on them as subcat.ID
				pc := []shared.PartComponents{}
				DB.SQL(`select * from component_part where component_id=$1`, c.ID).QueryStructs(&pc)
				for _, thePart := range pc {
					DB.SQL(`update part set category=$1 where id=$2`, subcat.ID, thePart.PartID).Exec()
				}

			}
		}

		*result = r
	}

	logger(start, "Util.Cats",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		*result,
		channel, conn.UserID, "part", 0, true)

	return nil
}

func (u *UtilRPC) AddPhoto(data shared.PhotoRPCData, newID *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Generate the thumbnail
	theImage := data.Photo.Photo[23:]
	// print("The Image =", theImage[:80])
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(theImage))
	m, _, err := image.Decode(reader)
	if err != nil {
		println("Decode Error", err.Error())
		// log.Fatal(err)
	} else {

		// bounds := m.Bounds()

		// // Calculate a 16-bin histogram for m's red, green, blue and alpha components.
		// //
		// // An image's bounds do not necessarily start at (0, 0), so the two loops start
		// // at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
		// // likely to result in better memory access patterns than X first and Y second.
		// var histogram [16][4]int
		// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		// 		r, g, b, a := m.At(x, y).RGBA()
		// 		// A color's RGBA method returns values in the range [0, 65535].
		// 		// Shifting by 12 reduces this to the range [0, 15].
		// 		histogram[r>>12][0]++
		// 		histogram[g>>12][1]++
		// 		histogram[b>>12][2]++
		// 		histogram[a>>12][3]++
		// 	}
		// }

		// // Print the results.
		// fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
		// for i, x := range histogram {
		// 	fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
		// }

		// create the thumbnail and a preview
		var tb bytes.Buffer
		thumb := resize.Resize(64, 0, m, resize.Lanczos3)
		encoder := base64.NewEncoder(base64.StdEncoding, &tb)
		jpeg.Encode(encoder, thumb, &jpeg.Options{Quality: 95})
		data.Photo.Thumbnail = "data:image/jpeg;base64," + tb.String()

		var pb bytes.Buffer
		preview := resize.Resize(240, 0, m, resize.Lanczos3)
		encoder = base64.NewEncoder(base64.StdEncoding, &pb)
		jpeg.Encode(encoder, preview, &jpeg.Options{Quality: 95})
		data.Photo.Preview = "data:image/jpeg;base64," + pb.String()

		// out, err := os.Create(fmt.Sprintf("public/thumbs/%d.jpg", *newID))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer out.Close()

		// // write new image to file
		// jpeg.Encode(out, thumb, nil)
	}

	// Save the data, and get a new ID
	DB.InsertInto("phototest").
		Columns("name", "photo", "thumbnail", "preview").
		Record(data.Photo).
		Returning("id").
		QueryScalar(newID)

	logger(start, "Util.AddPhoto",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		data.Photo.Name,
		data.Channel, conn.UserID, "phototest", 0, true)

	return nil
}

func (u *UtilRPC) GetPhoto(data shared.PhotoRPCData, photo *shared.Photo) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select id,name,preview from phototest where id=$1`, data.ID).QueryStruct(photo)

	logger(start, "Util.GetPhoto",
		fmt.Sprintf("Channel %d, ID %d, User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		photo.Name,
		data.Channel, conn.UserID, "phototest", data.ID, false)

	return nil
}

func (u *UtilRPC) PhotoList(data shared.PhotoRPCData, photos *[]shared.Photo) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select id,name,thumbnail from phototest order by name`).QueryStructs(photos)

	logger(start, "Util.PhotoList",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Photos", len(*photos)),
		data.Channel, conn.UserID, "photo", 0, false)

	return nil
}
