package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type PartRPC struct{}

// Get the details for a given part
func (p *PartRPC) Get(partID int, part *shared.Part) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(`select * from part where id=$1`, partID).QueryStruct(part)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.Get",
		fmt.Sprintf("Part %d", partID),
		part.Name)

	return nil
}

// Get the details for a given part class
func (p *PartRPC) GetClass(id int, partClass *shared.PartClass) error {
	start := time.Now()

	if id == 0 {
		*partClass = shared.PartClass{
			ID:    0,
			Name:  "All",
			Descr: "Parts that apply to all machine types",
		}
	} else {
		// Read the sites that this user has access to
		err := DB.SQL(`select * from part_class where id=$1`, id).QueryStruct(partClass)

		if err != nil {
			log.Println(err.Error())
		}
	}

	logger(start, "Part.GetClass",
		fmt.Sprintf("Class %d", id),
		partClass.Name)

	return nil
}

// Get a list of machine classes
func (m *PartRPC) ClassList(channel int, classes *[]shared.PartClass) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Read the sites that this user has access to
	*classes = append(*classes, shared.PartClass{
		ID:    0,
		Name:  "All",
		Descr: "Parts that apply to all machine types",
	})
	err := DB.SQL(`select id,name,descr from part_class order by name`).
		QueryStructs(classes)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.ClassList",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Classes", len(*classes)))

	return nil
}

// Get all the parts for the given class
func (p *PartRPC) List(req shared.PartListReq, parts *[]shared.Part) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select * from part 
		where class=$1
		order by name`, req.Class).QueryStructs(parts)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			req.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("Class %s %d parts", req.Class, len(*parts)))

	return nil
}
