package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type MachineRPC struct{}

const MachineQuery = `select 
m.*,s.name as site_name,x.span as span
		from machine m
		left join site s on (s.id=m.site_id)
		left join site_layout x on (x.site_id=m.site_id and x.machine_id=m.id)
where m.id = $1`

// Get the details for a given machine
func (m *MachineRPC) Get(machineID int, machine *shared.Machine) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(MachineQuery, machineID).QueryStruct(machine)

	if err != nil {
		log.Println(err.Error())
	}

	// fetch all components
	err = DB.Select("*").
		From("component").
		Where("machine_id = $1", machineID).
		OrderBy("position,zindex,lower(name)").
		QueryStructs(&machine.Components)

	logger(start, "Machine.Get",
		fmt.Sprintf("Machine %d", machineID),
		machine.Name)

	return nil
}

// Update a Machine
func (m *MachineRPC) Update(data shared.MachineUpdateData, ok *bool) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.Update("machine").
		SetWhitelist(data.Machine, "name", "serialnum", "descr", "notes").
		Where("id = $1", data.Machine.ID).
		Exec()

	logger(start, "Machine.Update",
		fmt.Sprintf("Channel %d, Machine %d, User %d %s %s",
			data.Channel, data.Machine.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Machine.Name)

	*ok = true
	return nil
}

// Insert a machine
func (m *MachineRPC) Insert(data shared.MachineUpdateData, id *int) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	*id = 0
	DB.InsertInto("machine").
		Columns("name", "serialnum", "descr", "notes", "site_id").
		Record(data.Machine).
		Returning("id").
		QueryScalar(id)

	logger(start, "Machine.Insert",
		fmt.Sprintf("Channel %d, Machine %d, User %d %s %s",
			data.Channel, id, conn.UserID, conn.Username, conn.UserRole),
		data.Machine.Name)

	return nil
}

// Delete a machine
func (m *MachineRPC) Delete(data shared.MachineUpdateData, ok *bool) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	*ok = false
	id := data.Machine.ID
	DB.DeleteFrom("machine").
		Where("id=$1", id).
		Exec()

	logger(start, "Machine.Delete",
		fmt.Sprintf("Channel %d, Machine %d, User %d %s %s",
			data.Channel, id, conn.UserID, conn.Username, conn.UserRole),
		data.Machine.Name)

	*ok = true
	return nil
}
