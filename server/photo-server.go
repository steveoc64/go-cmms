package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"strings"
	"time"

	"itrak-cmms/shared"

	"github.com/nfnt/resize"
)

func decodePhoto(photo string, preview *string, thumbnail *string) error {

	if photo == "" || len(photo) < 22 {
		print("photo is empty")
		*preview = ""
		*thumbnail = ""
		return nil
	}
	theImage := ""

	// println("passed in", photo)
	// println("Decode Photo Data =", photo[:80], "...")
	f := strings.SplitN(photo, ",", 2)
	switch f[0] {
	case "data:image/jpeg;base64":
		theImage = f[1]
	case "data:image/png;base64":
		theImage = f[1]
	case "data:application/pdf;base64":
		*preview = PDFPreview
		*thumbnail = PDFThumb
		return nil
	default:
		println("unknown file format")
		return nil
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(theImage))
	m, _, err := image.Decode(reader)
	if err != nil {
		println("Decode Error", err.Error())
		// log.Fatal(err)
	} else {
		// create the thumbnail and a preview
		var tb bytes.Buffer
		thumbVar := resize.Resize(64, 0, m, resize.Lanczos3)
		encoder := base64.NewEncoder(base64.StdEncoding, &tb)
		jpeg.Encode(encoder, thumbVar, &jpeg.Options{Quality: 95})
		*thumbnail = "data:image/jpeg;base64," + tb.String()

		var pb bytes.Buffer
		previewVar := resize.Resize(240, 0, m, resize.Lanczos3)
		encoder = base64.NewEncoder(base64.StdEncoding, &pb)
		jpeg.Encode(encoder, previewVar, &jpeg.Options{Quality: 95})
		*preview = "data:image/jpeg;base64," + pb.String()
	}

	return nil
}

func (u *UtilRPC) AddPhoto(data shared.PhotoRPCData, newID *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	// print("addphoto", data)
	// print("addphoto", data.Photo)
	// print("addphoto", data.Photo.Photo)

	decodePhoto(data.Photo.Photo, &data.Photo.Preview, &data.Photo.Thumb)

	// Save the data, and get a new ID
	DB.InsertInto("photo").
		Columns("notes", "photo", "thumb", "preview").
		Record(data.Photo).
		Returning("id").
		QueryScalar(newID)

	DB.SQL(`update photo set entity='test', entity_id=$1 where id=$1`, newID).Exec()

	logger(start, "Util.AddPhoto",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d", *newID),
		data.Channel, conn.UserID, "photo", 0, true)

	return nil
}

func (u *UtilRPC) GetPhoto(data shared.PhotoRPCData, photo *shared.Photo) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select id,notes,preview,entity,entity_id from photo where id=$1`, data.ID).QueryStruct(photo)

	logger(start, "Util.GetPhoto",
		fmt.Sprintf("Channel %d, ID %d, User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		photo.Notes,
		data.Channel, conn.UserID, "photo", data.ID, false)

	return nil
}

func (u *UtilRPC) GetFullPhoto(data shared.PhotoRPCData, photo *shared.Photo) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select id,notes,photo,preview,entity,entity_id from photo where id=$1`, data.ID).QueryStruct(photo)

	logger(start, "Util.GetFullPhoto",
		fmt.Sprintf("Channel %d, ID %d, User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		photo.Photo[:22],
		data.Channel, conn.UserID, "photo", data.ID, false)

	return nil
}

func (u *UtilRPC) PhotoList(data shared.PhotoRPCData, photos *[]shared.Photo) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select id,entity,entity_id,notes,thumb from photo order by id desc`).QueryStructs(photos)

	logger(start, "Util.PhotoList",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Photos", len(*photos)),
		data.Channel, conn.UserID, "photo", 0, false)

	return nil
}

func (u *UtilRPC) UpdatePhoto(data shared.PhotoRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// decodePhoto(data.Photo.Photo, &data.Photo.Preview, &data.Photo.Thumb)

	// Save the data
	DB.Update("photo").
		SetWhitelist(data.Photo, "notes", "entity", "entity_id").
		Where("id = $1", data.ID).
		Exec()

	logger(start, "Util.UpdatePhoto",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d", data.ID),
		data.Channel, conn.UserID, "photo", 0, true)

	*done = true

	return nil
}

func (u *UtilRPC) DeletePhoto(data shared.PhotoRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Save the data
	DB.SQL(`delete from photo where id=$1`, data.ID).Exec()

	logger(start, "Util.DeletePhoto",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d", data.ID),
		data.Channel, conn.UserID, "photo", 0, true)

	*done = true

	return nil
}

var PDFImage string
var PDFPreview string
var PDFThumb string

func cachePDFImage() {
	id := 0
	PDFImage = ""
	DB.SQL(`select id,photo,preview,thumb from stdimg where code='PDF'`).QueryScalar(&id, &PDFImage, &PDFPreview, &PDFThumb)
	if id > 0 {
		fmt.Printf("Cached PDF Image %d len %d header %s\n", id, len(PDFImage), PDFImage[:44])
	} else {
		println("*** No standard PDF Image in database ... please fix !!! ***")
	}
}

func (u *UtilRPC) GetPDF(channel int, pdf *string) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Save the data
	id := 0
	*pdf = ""
	DB.SQL(`select id,photo from stdimg where code='PDF'`).QueryScalar(&id, pdf)

	logger(start, "Util.GetPDF",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Img ID %d size %d header %s", id, len(*pdf), (*pdf)[:44]),
		channel, conn.UserID, "stdimg", id, false)

	return nil
}
