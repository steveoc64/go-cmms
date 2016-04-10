package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type PartRPC struct{}

// Get the details for a given machine
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

// Get all the parts
func (p *PartRPC) List(channel int, parts *[]shared.Part) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(`select * from part order by name`).QueryStructs(parts)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.List",
		fmt.Sprintf("Channel %d", channel),
		fmt.Sprintf("%d parts", len(*parts)))

	return nil
}
