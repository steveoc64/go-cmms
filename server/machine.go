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
func (m *MachineRPC) Get(data shared.MachineRPCData, machine *shared.Machine) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the sites that this user has access to
	err := DB.SQL(MachineQuery, data.ID).QueryStruct(machine)

	if err != nil {
		log.Println(err.Error())
	}

	// fetch all components
	err = DB.Select("*").
		From("component").
		Where("machine_id = $1", data.ID).
		OrderBy("position,zindex,lower(name)").
		QueryStructs(&machine.Components)

	logger(start, "Machine.Get",
		fmt.Sprintf("%d", data.ID),
		machine.Name,
		data.Channel, conn.UserID, "machine", data.ID, false)

	return nil
}

// Update a Machine
func (m *MachineRPC) Update(data shared.MachineRPCData, ok *bool) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.Update("machine").
		SetWhitelist(data.Machine, "name", "serialnum", "descr", "notes",
			"alerts_to", "tasks_to", "part_class").
		Where("id = $1", data.Machine.ID).
		Exec()

	logger(start, "Machine.Update",
		fmt.Sprintf("Channel %d, Machine %d, User %d %s %s",
			data.Channel, data.Machine.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Machine.Name,
		data.Channel, conn.UserID, "machine", data.Machine.ID, true)

	*ok = true
	return nil
}

// Insert a machine
func (m *MachineRPC) Insert(data shared.MachineRPCData, id *int) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	*id = 0
	DB.InsertInto("machine").
		Columns("name", "serialnum", "descr", "notes", "site_id",
			"alerts_to", "tasks_to", "part_class").
		Record(data.Machine).
		Returning("id").
		QueryScalar(id)

	logger(start, "Machine.Insert",
		fmt.Sprintf("Channel %d, Machine %d, User %d %s %s",
			data.Channel, *id, conn.UserID, conn.Username, conn.UserRole),
		data.Machine.Name,
		data.Channel, conn.UserID, "machine", *id, true)

	return nil
}

// Delete a machine
func (m *MachineRPC) Delete(data shared.MachineRPCData, ok *bool) error {
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
		data.Machine.Name,
		data.Channel, conn.UserID, "machine", id, true)

	*ok = true
	return nil
}

func (m *MachineRPC) MachineTypes(data shared.MachineRPCData, machineTypes *[]shared.MachineType) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.Select(`*`).From(`machine_type`).OrderBy(`name`).QueryStructs(machineTypes)

	logger(start, "Machine.MachineTypes",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d machine types", len(*machineTypes)),
		data.Channel, conn.UserID, "machine_type", 0, true)

	return nil
}

func (m *MachineRPC) GetMachineType(data shared.MachineTypeRPCData, machineType *shared.MachineType) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.Select(`*`).
		From(`machine_type`).
		Where(`id=$1`, data.ID).
		QueryStruct(machineType)

	logger(start, "Machine.GetMachineType",
		fmt.Sprintf("Channel %d, ID %d User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		machineType.Name,
		data.Channel, conn.UserID, "machine_type", data.ID, true)

	return nil
}

func (m *MachineRPC) UpdateMachineType(data shared.MachineTypeRPCData, done *bool) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.SQL(`update machine_type set name=$2 where id=$1`, data.ID, data.MachineType.Name).Exec()

	logger(start, "Machine.UpdateMachineType",
		fmt.Sprintf("Channel %d, ID %d User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		data.MachineType.Name,
		data.Channel, conn.UserID, "machine_type", data.ID, true)

	*done = true
	return nil
}
