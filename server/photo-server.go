package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	_ "image/png"
	"strings"
	"time"

	"itrak-cmms/shared"

	"github.com/nfnt/resize"
)

//func decodePhoto(photo string, preview *string, thumbnail *string, thetype *string, datatype *string) error {
func decodePhoto(photo *shared.Photo) error {

	if photo.Data == "" || len(photo.Data) < 22 {
		print("photo is empty")
		photo.Preview = ""
		photo.Thumb = ""
		return nil
	}
	theImage := ""

	// println("passed in", photo)
	// println("Decode Photo Data =", photo[:80], "...")
	f := strings.SplitN(photo.Data, ",", 2)
	photo.Datatype = f[0]
	switch f[0] {
	case "data:image/jpeg;base64", "data:image/png;base64", "data:image/gif;base64":
		theImage = f[1]
		photo.Type = "Image"
	case "data:application/pdf;base64":
		photo.Preview = PDFPreview
		photo.Thumb = PDFThumb
		photo.Type = "PDF"
		return nil
	default:
		println("unknown file format", f[0])
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
		// jpeg.Encode(encoder, thumbVar, &jpeg.Options{Quality: 50})
		gif.Encode(encoder, thumbVar, &gif.Options{NumColors: 256})
		// photo.Thumb = "data:image/jpeg;base64," + tb.String()
		photo.Thumb = "data:image/gif;base64," + tb.String()

		var pb bytes.Buffer
		previewVar := resize.Resize(240, 0, m, resize.Lanczos3)
		encoder = base64.NewEncoder(base64.StdEncoding, &pb)
		// jpeg.Encode(encoder, previewVar, &jpeg.Options{Quality: 75})
		gif.Encode(encoder, previewVar, &gif.Options{NumColors: 256})
		// photo.Preview = "data:image/jpeg;base64," + pb.String()
		photo.Preview = "data:image/gif;base64," + pb.String()
	}

	return nil
}

func (u *UtilRPC) AddPhoto(data shared.PhotoRPCData, newID *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	// print("addphoto", data)
	// print("addphoto", data.Photo)
	// print("addphoto", data.Photo.Photo)

	//decodePhoto(data.Photo.Data, &data.Photo.Preview, &data.Photo.Thumb, &data.Photo.Type, &data.Photo.Datatype)
	decodePhoto(data.Photo)

	// Save the data, and get a new ID
	DB.InsertInto("photo").
		Columns("notes", "photo", "thumb", "preview", "filename", "type", "datatype").
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

	DB.SQL(`select id,notes,preview,entity,entity_id,filename from photo where id=$1`, data.ID).QueryStruct(photo)

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

	DB.SQL(`select id,notes,photo,preview,entity,entity_id,filename from photo where id=$1`, data.ID).QueryStruct(photo)

	logger(start, "Util.GetFullPhoto",
		fmt.Sprintf("Channel %d, ID %d, User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		photo.Data[:22],
		data.Channel, conn.UserID, "photo", data.ID, false)

	return nil
}

func (u *UtilRPC) PhotoList(data shared.PhotoTestRPCData, photos *[]shared.Photo) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select id,entity,entity_id,notes,thumb,filename from photo order by id desc`).QueryStructs(photos)

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
