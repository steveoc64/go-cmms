package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

// PartRPC exported struct for catching RPC calls
type PartRPC struct{}

// Get the details for a given part
func (p *PartRPC) Get(data shared.PartRPCData, part *shared.Part) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select * from part where id=$1`, data.ID).QueryStruct(part)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.Get",
		fmt.Sprintf("%d", data.ID),
		part.Name,
		data.Channel, conn.UserID, "part", data.ID, false)

	return nil
}

// GetClass - Get the details for a given part class
func (p *PartRPC) GetClass(data shared.PartClassRPCData, partClass *shared.PartClass) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	id := data.ID

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
		fmt.Sprintf("%d", id),
		partClass.Name,
		data.Channel, conn.UserID, "part_class", id, false)

	return nil
}

// InsertClass - Add a new part class
func (p *PartRPC) InsertClass(data shared.PartClassRPCData, id *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	*id = 0
	DB.InsertInto("part_class").
		Columns("name", "descr").
		Record(data.PartClass).
		Returning("id").
		QueryScalar(id)

	logger(start, "Part.InsertClass",
		data.PartClass.Name,
		"",
		data.Channel, conn.UserID, "part_class", *id, true)

	return nil
}

// DeleteClass - Delete the class
func (p *PartRPC) DeleteClass(data shared.PartClassRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.DeleteFrom("part_class").
		Where("id=$1", data.PartClass.ID).
		Exec()

	logger(start, "Part.DeleteClass",
		fmt.Sprintf("Channel %d, Class %d, User %d %s %s",
			data.Channel, data.PartClass.ID, conn.UserID, conn.Username, conn.UserRole),
		data.PartClass.Name,
		data.Channel, conn.UserID, "part_class", data.PartClass.ID, true)

	*done = true

	return nil
}

// UpdateClass - Update the class
func (p *PartRPC) UpdateClass(data shared.PartClassRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("part_class").
		SetWhitelist(data.PartClass, "name", "descr").
		Where("id = $1", data.PartClass.ID).
		Exec()

	logger(start, "Part.UpdateClass",
		fmt.Sprintf("%s : %s", data.PartClass.Name, data.PartClass.Descr),
		"",
		data.Channel, conn.UserID, "part_class", data.PartClass.ID, true)

	*done = true

	return nil
}

// ClassList - Get a list of machine classes
func (p *PartRPC) ClassList(channel int, classes *[]shared.PartClass) error {
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
		"",
		fmt.Sprintf("%d Classes", len(*classes)),
		channel, conn.UserID, "part_class", 0, false)

	return nil
}

// List - Get all the parts for the given class, which is passed in as the ID
// or leave the ID 0 to get all parts
func (p *PartRPC) List(data shared.PartRPCData, parts *[]shared.Part) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	if data.ID == 0 {
		DB.SQL(`select * from part order by name`).QueryStructs(parts)
	} else {
		DB.SQL(`select * from part
			where class=$1
			order by name`, data.ID).
			QueryStructs(parts)
	}

	logger(start, "Part.List",
		"",
		fmt.Sprintf("Class %d %d parts", data.ID, len(*parts)),
		data.Channel, conn.UserID, "parts", data.ID, false)

	return nil
}

// Update the part
func (p *PartRPC) Update(data shared.PartRPCData, done *bool) error {
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
		data.Part.Name,
		"",
		data.Channel, conn.UserID, "part", data.Part.ID, true)

	return nil
}

// Insert a new part
func (p *PartRPC) Insert(data shared.PartRPCData, id *int) error {
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
		data.Part.Name,
		fmt.Sprintf("%d", *id),
		data.Channel, conn.UserID, "part", *id, true)

	return nil
}

// Delete a new part
func (p *PartRPC) Delete(data shared.PartRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.DeleteFrom("part").
		Where("id=$1", data.Part.ID).
		Exec()

	DB.DeleteFrom("part_price").
		Where("part_id=$1", data.Part.ID).
		Exec()

	logger(start, "Part.Delete",
		fmt.Sprintf("%d", data.Part.ID),
		data.Part.Name,
		data.Channel, conn.UserID, "part", data.Part.ID, true)

	*done = true
	return nil
}

// StockList - Get a list of stock records for a part
func (p *PartRPC) StockList(data shared.PartRPCData, stocks *[]shared.PartStock) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the stock records for this part in reverse date order
	err := DB.SQL(`select * 
		from part_stock 
		where part_id=$1 
		order by datefrom desc
		limit 5`, data.ID).
		QueryStructs(stocks)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.StockList",
		fmt.Sprintf("%d", data.ID),
		fmt.Sprintf("%d stock records", len(*stocks)),
		data.Channel, conn.UserID, "part_stock", 0, false)

	return nil
}

// PriceList - Get a list of price records for a part
func (p *PartRPC) PriceList(data shared.PartRPCData, prices *[]shared.PartPrice) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the stock records for this part in reverse date order
	err := DB.SQL(`select * 
		from part_price 
		where part_id=$1 
		order by datefrom desc
		limit 5`, data.ID).
		QueryStructs(prices)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Part.PriceList",
		fmt.Sprintf("%d", data.ID),
		fmt.Sprintf("%d price records", len(*prices)),
		data.Channel, conn.UserID, "part_price", 0, false)

	return nil
}

// GetCategory - Get the category record
func (p *PartRPC) GetCategory(data shared.PartRPCData, cat *shared.Category) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)
	DB.SQL(`select * from category where id=$1`, data.ID).QueryStruct(cat)

	logger(start, "Part.GetCategory",
		fmt.Sprintf("%d", data.ID),
		fmt.Sprintf("%v", cat),
		data.Channel, conn.UserID, "category", data.ID, false)

	return nil
}

// AddCategory - Add a category with the specified Cat as the parent, return the new Cat ID
func (p *PartRPC) AddCategory(data shared.PartRPCData, newID *int) error {
	start := time.Now()

	// Calc the stock code based on the parent, with the seq number attached
	stockCode := ""
	if data.ID != 0 {
		pcat := shared.Category{}
		numSubcats := 0
		DB.SQL(`select stock_code from category where id=$1`, data.ID).QueryStruct(&pcat)
		DB.SQL(`select count(*) from category where parent_id=$1`, data.ID).QueryScalar(&numSubcats)
		stockCode = fmt.Sprintf("%s-%02d", pcat.StockCode, numSubcats+1)
		print("new stock code =", stockCode)
	}

	conn := Connections.Get(data.Channel)
	cat := shared.Category{
		ParentID:  data.ID,
		Name:      "New Category",
		StockCode: stockCode,
	}
	DB.InsertInto("category").
		Columns("parent_id", "name", "stock_code").
		Record(&cat).
		Returning("id").
		QueryScalar(newID)

	logger(start, "Part.AddCategory",
		fmt.Sprintf("Channel %d Parent %d", data.Channel, data.ID),
		fmt.Sprintf("New ID %d", *newID),
		data.Channel, conn.UserID, "category", *newID, true)

	return nil
}

// AddPart - Add a part with the specified Cat as the parent, return the new part ID
func (p *PartRPC) AddPart(data shared.PartRPCData, newPart *shared.Part) error {
	start := time.Now()

	stockCode := ""
	DB.SQL(`select stock_code from category where id=$1`, data.ID).QueryScalar(&stockCode)

	conn := Connections.Get(data.Channel)
	part := shared.Part{
		Category:  data.ID,
		Name:      "New Part",
		StockCode: fmt.Sprintf("%s-0000", stockCode),
	}
	newID := 0
	DB.InsertInto("part").
		Columns("category", "name", "stock_code").
		Record(&part).
		Returning("id").
		QueryScalar(&newID)

	DB.SQL(`select * from part where id=$1`, newID).QueryStruct(newPart)

	logger(start, "Part.AddPart",
		fmt.Sprintf("Channel %d Category %d", data.Channel, data.ID),
		fmt.Sprintf("New ID %d", newPart.ID),
		data.Channel, conn.UserID, "part", newPart.ID, true)

	return nil
}

// DelPart - Delete the part in data.ID, and return the parent category ID
func (p *PartRPC) DelPart(data shared.PartRPCData, parentCat *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.SQL(`select category from part where id=$1`, data.ID).QueryScalar(parentCat)

	DB.DeleteFrom("part").
		Where("id=$1", data.ID).
		Exec()

	logger(start, "Part.DelPart",
		fmt.Sprintf("Channel %d Part %d", data.Channel, data.ID),
		fmt.Sprintf("Parent Cat %d", *parentCat),
		data.Channel, conn.UserID, "part", data.ID, true)

	return nil
}

// GetTree - Get a parts tree from a specifec category ... uses recursive Fn getTree() to complete
func (p *PartRPC) GetTree(data shared.PartTreeRPCData, cats *[]shared.Category) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	c := getTree(data.CategoryID)
	*cats = c

	logger(start, "Part.GetTree",
		fmt.Sprintf("%d", data.CategoryID),
		fmt.Sprintf("%d subcats", len(*cats)),
		data.Channel, conn.UserID, "category", 0, false)

	return nil
}

// getTree - Recursive Fn to fill in the nodes on a parts tree
func getTree(parentCat int) []shared.Category {
	// fmt.Printf("getting categories with parent %d\n", parentCat)

	cats := []shared.Category{}

	// get an array of categories that point to the given parent
	DB.SQL(`select * from category where parent_id=$1 order by name`, parentCat).QueryStructs(&cats)
	// fmt.Printf("subcats of %d = %v (%d)\n", parentCat, cats, len(cats))

	for i, c := range cats {
		DB.SQL(`select * from part where category=$1 order by name`, c.ID).QueryStructs(&cats[i].Parts)
		// fmt.Printf("parts of cat %d = %v (%d)\n", c.ID, cats[i].Parts, len(cats[i].Parts))
		cats[i].Subcats = getTree(c.ID)
	}

	return cats
}
