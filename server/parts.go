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

// Add a new part class
func (p *PartRPC) InsertClass(data shared.PartClassUpdateData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	*id = 0
	DB.InsertInto("part_class").
		Columns("name", "descr").
		Record(data.PartClass).
		Returning("id").
		QueryScalar(id)

	logger(start, "Part.InsertClass",
		fmt.Sprintf("Channel %d, Class %d, User %d %s %s",
			data.Channel, *id, conn.UserID, conn.Username, conn.UserRole),
		data.PartClass.Name)

	return nil
}

// Delete the class
func (p *PartRPC) DeleteClass(data shared.PartClassUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.DeleteFrom("part_class").
		Where("id=$1", data.PartClass.ID).
		Exec()

	logger(start, "Part.DeleteClass",
		fmt.Sprintf("Channel %d, Class %d, User %d %s %s",
			data.Channel, data.PartClass.ID, conn.UserID, conn.Username, conn.UserRole),
		data.PartClass.Name)

	*done = true

	return nil
}

// Update the class
func (p *PartRPC) UpdateClass(data shared.PartClassUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("part_class").
		SetWhitelist(data.PartClass, "name", "descr").
		Where("id = $1", data.PartClass.ID).
		Exec()

	logger(start, "Part.UpdateClass",
		fmt.Sprintf("Channel %d, Class %d, User %d %s %s",
			data.Channel, data.PartClass.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%s : %s", data.PartClass.Name, data.PartClass.Descr))

	*done = true

	return nil
}

// Get a list of machine classes
func (m *PartRPC) ClassList(channel int, classes *[]shared.PartClass) error {
	start := time.Now()

	conn := Connections.Get(channel)

	haveNone := 0
	DB.SQL(`select count(*) from part where part.class=0`).QueryScalar(&haveNone)

	// Read the sites that this user has access to
	*classes = append(*classes, shared.PartClass{
		ID:    0,
		Name:  "All",
		Descr: "Parts that apply to all machine types",
		Count: haveNone,
	})

	err := DB.SQL(`select 
		p.id as id,p.name as name,p.descr as descr,
		(select count(*) from part where part.class=p.id) as count
		from part_class p order by p.name`).
		QueryStructs(classes)

	// err := DB.SQL(`select id,name,descr from part_class order by name`).
	// 	QueryStructs(classes)

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
		fmt.Sprintf("Class %d %d parts", req.Class, len(*parts)))

	return nil
}

// Update the part
func (p *PartRPC) Update(data shared.PartUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// get the last price and stock level
	existingPart := shared.Part{}
	DB.SQL(`select * from part where id=$1`, data.Part.ID).QueryStruct(&existingPart)

	DB.Update("part").
		SetWhitelist(data.Part,
			"class", "name", "descr", "stock_code", "reorder_stocklevel",
			"reorder_qty", "latest_price", "qty_type", "notes", "current_stock").
		Where("id = $1", data.Part.ID).
		Exec()

	*done = true

	if existingPart.CurrentStock != data.Part.CurrentStock {
		// create a new part_stock record
		partStock := shared.PartStock{
			PartID:     data.Part.ID,
			StockLevel: data.Part.CurrentStock,
			Descr:      fmt.Sprintf("Updated by %s", conn.Username),
		}
		DB.InsertInto("part_stock").
			Columns("part_id", "stock_level").
			Record(partStock).
			Exec()
		*done = false
	}

	if existingPart.LatestPrice != data.Part.LatestPrice {
		// update the last price date, and create a new part_price record
		DB.SQL(`update part set last_price_date=now() where id=$1`,
			data.Part.ID,
			fmt.Sprintf("Updated by %s", conn.Username)).Exec()

		partPrice := shared.PartPrice{
			PartID: data.Part.ID,
			Price:  data.Part.LatestPrice,
			Descr:  fmt.Sprintf("Updated by %s", conn.Username),
		}
		DB.InsertInto("part_price").
			Columns("part_id", "price", "descr").
			Record(partPrice).
			Exec()
		*done = false
	}

	logger(start, "Part.Update",
		fmt.Sprintf("Channel %d, Part %d, User %d %s %s",
			data.Channel, data.Part.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Part.Name)

	return nil
}

// Insert a new part
func (p *PartRPC) Insert(data shared.PartUpdateData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.InsertInto("part").
		Columns("class", "name", "descr", "stock_code", "reorder_stocklevel",
			"reorder_qty", "latest_price", "qty_type", "notes", "current_stock").
		Record(data.Part).
		Returning("id").
		QueryScalar(id)

	// create a new part_stock record
	partStock := shared.PartStock{
		PartID:     *id,
		StockLevel: data.Part.CurrentStock,
	}
	DB.InsertInto("part_stock").
		Columns("part_id", "stock_level").
		Record(partStock).
		Exec()

	// update the last price date, and create a new part_price record
	DB.SQL(`update part set last_price_date=now(), where id=$1`, *id).Exec()

	partPrice := shared.PartPrice{
		PartID: *id,
		Price:  data.Part.LatestPrice,
	}
	DB.InsertInto("part_price").
		Columns("part_id", "price").
		Record(partPrice).
		Exec()

	logger(start, "Part.Insert",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("New Part %d", *id))

	return nil
}

// Delete a new part
func (p *PartRPC) Delete(data shared.PartUpdateData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.DeleteFrom("part").
		Where("id=$1", data.Part.ID).
		Exec()

	DB.DeleteFrom("part_price").
		Where("part_id=$1", data.Part.ID).
		Exec()

	logger(start, "Part.Delete",
		fmt.Sprintf("Channel %d, Part %d, User %d %s %s",
			data.Channel, data.Part.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Part.Name)

	*done = true
	return nil
}

// Get a list of stock records for a part
func (p *PartRPC) StockList(id int, stocks *[]shared.PartStock) error {
	start := time.Now()

	// Read the stock records for this part in reverse date order
	err := DB.SQL(`select * 
		from part_stock 
		where part_id=$1 
		order by datefrom desc
		limit 5`, id).
		QueryStructs(stocks)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.StockList",
		fmt.Sprintf("Part %d", id),
		fmt.Sprintf("%d stock records", len(*stocks)))

	return nil
}

// Get a list of price records for a part
func (p *PartRPC) PriceList(id int, prices *[]shared.PartPrice) error {
	start := time.Now()

	// Read the stock records for this part in reverse date order
	err := DB.SQL(`select * 
		from part_price 
		where part_id=$1 
		order by datefrom desc
		limit 5`, id).
		QueryStructs(prices)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.PriceList",
		fmt.Sprintf("Part %d", id),
		fmt.Sprintf("%d price records", len(*prices)))

	return nil
}
