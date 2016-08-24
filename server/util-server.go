package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os/exec"
	"strings"
	"time"

	"itrak-cmms/shared"
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

// Patch up the Machine Components to point to the correct Machine Type Tool by ID
func (u *UtilRPC) MTT(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)
	*result = ""

	if conn.UserRole == "Admin" && conn.Username == "steve" {
		r := "Processing Machine Type Tools\n"

		components := []shared.Component{}

		DB.SQL(`select * from component order by machine_id,position`).QueryStructs(&components)

		mt := 0
		mtt := 0

		// patched := 0
		for _, c := range components {

			DB.SQL(`select machine_type from machine where id=$1`, c.MachineID).QueryScalar(&mt)
			DB.SQL(`select id from machine_type_tool where machine_id=$1 and position=$2`, mt, c.Position).QueryScalar(&mtt)
			r += fmt.Sprintf("Component ID %d: Machine %d:%d MT %d MTT %d  %s\n",
				c.ID,
				c.MachineID, c.Position,
				mt, mtt,
				c.Name)
			DB.SQL(`update component set mtt_id=$1 where id=$2`, mtt, c.ID).Exec()

		}
		*result = r

	}

	logger(start, "Util.MTT",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		*result,
		channel, conn.UserID, "component", 0, true)

	return nil
}

// Move all photos into their own table
func (u *UtilRPC) PhotoMove(channel int, result *string) error {
	start := time.Now()

	conn := Connections.Get(channel)
	*result = ""

	if conn.UserRole == "Admin" && conn.Username == "steve" {
		r := "Processing Photos into their own tables\n"

		println("events")
		events := []shared.Event{}
		err := DB.SQL(`select id,photo,photo_preview,photo_thumbnail from event where length(photo)>0`).QueryStructs(&events)

		// patched := 0
		for k, v := range events {
			r += fmt.Sprintf("Event %d\n", v.ID)
			println(k, ":", v.ID)

			id := 0
			pf := strings.Split(v.Photo.Data, `,`)
			phototype := pf[0]

			if v.Photo.Data != "" {
				err = DB.SQL(`insert into photo (datatype,entity,entity_id,photo,preview,thumb) values($1,$2,$3,$4,$5,$6) returning id`,
					phototype, `event`, v.ID, v.Photo.Data, v.Photo.Preview, v.Photo.Thumb).
					QueryScalar(&id)
				if err != nil {
					r += err.Error()
				}
				r += fmt.Sprintf("Event Photo1 %d = %d\n", v.ID, id)
			}
		}

		// now strip the photos from the events table
		// DB.SQL(`alter table event drop photo`).Exec()
		// DB.SQL(`alter table event drop photo_preview`).Exec()
		// DB.SQL(`alter table event drop photo_thumbnail`).Exec()

		println("tasks")
		tasks := []shared.Task{}
		DB.SQL(`select id,photo1,photo2,photo3,preview1,preview2,preview3,thumb1,thumb2,thumb3 from task where length(photo1)>0`).QueryStructs(&tasks)

		// patched := 0
		for k, v := range tasks {
			r += fmt.Sprintf("Task %d\n", v.ID)
			println(k, ":", v.ID)

			id := 0
			pf := strings.Split(v.Photo1, `,`)
			phototype := pf[0]

			if v.Photo1 != "" {
				err = DB.SQL(`insert into photo (datatype,entity,entity_id,photo,preview,thumb) values($1,$2,$3,$4,$5,$6) returning id`,
					phototype, `task`, v.ID, v.Photo1, v.Preview1, v.Thumb1).
					QueryScalar(&id)
				println("Added photo", id)
				r += fmt.Sprintf("Task Photo1 %d = %d\n", v.ID, id)
			}

			if v.Photo2 != "" {
				err = DB.SQL(`insert into photo (datatype,entity,entity_id,photo,preview,thumb) values($1,$2,$3,$4,$5,$6) returning id`,
					phototype, `task`, v.ID, v.Photo2, v.Preview2, v.Thumb2).
					QueryScalar(&id)
				r += fmt.Sprintf("Task Photo2 %d = %d\n", v.ID, id)
			}

			if v.Photo3 != "" {
				err = DB.SQL(`insert into photo (datatype,entity,entity_id,photo,preview,thumb) values($1,$2,$3,$4,$5,$6) returning id`,
					phototype, `task`, v.ID, v.Photo3, v.Preview3, v.Thumb3).
					QueryScalar(&id)
				r += fmt.Sprintf("Task Photo3 %d = %d\n", v.ID, id)
			}
			if err != nil {
				r += err.Error()
			}
		}

		// Copy over the phototest elements
		ptest := []shared.Phototest{}
		DB.SQL(`select * from phototest`).QueryStructs(&ptest)
		for k, v := range ptest {
			r += fmt.Sprintf("PhotoTest %d\n", v.ID)
			println(k, ":", v.ID)

			id := 0
			pf := strings.Split(v.Photo.Data, `,`)
			phototype := pf[0]

			if v.Photo.Data != "" {
				DB.SQL(`insert into photo (datatype,entity,entity_id,photo,preview,thumb,filename) values($1,$2,$3,$4,$5,$6) returning id`,
					phototype, `test`, v.ID, v.Photo, v.Preview, v.Thumbnail).
					QueryScalar(&id)
				println("Added photo", id)
				r += fmt.Sprintf("Test Photo %d = %d\n", v.ID, id)
			}
		}

		*result = r

	}

	logger(start, "Util.PhotoMove",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		*result,
		channel, conn.UserID, "photo", 0, true)

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
