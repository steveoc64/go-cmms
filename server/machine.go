package main

import (
	"fmt"
	"log"
	"time"

	"itrak-cmms/shared"
)

type MachineRPC struct{}

const MachineQuery = `select 
m.*,s.name as site_name,x.span as span
		from machine m
		left join site s on (s.id=m.site_id)
		left join site_layout x on (x.site_id=m.site_id and x.machine_id=m.id)
where m.id = $1`

const MachinesOfType = `select 
m.*,s.name as site_name,x.span as span
		from machine m
		left join site s on (s.id=m.site_id)
		left join site_layout x on (x.site_id=m.site_id and x.machine_id=m.id)
where m.machine_type = $1
order by x.seq,lower(m.name)`

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

// Get machines of a specific type
func (m *MachineRPC) MachineOfType(data shared.MachineRPCData, machines *[]shared.Machine) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	// Read the machines for the given site
	err := DB.SQL(MachinesOfType, data.ID).QueryStructs(machines)

	if err != nil {
		log.Println(err.Error())
	}

	// For each machine, fetch all components
	for k, m := range *machines {
		err = DB.Select("*").
			From("component").
			Where("machine_id = $1", m.ID).
			OrderBy("position,zindex,lower(name)").
			QueryStructs(&(*machines)[k].Components)
	}

	logger(start, "Machine.MachinesOfType",
		fmt.Sprintf("Channel %d, Type %d, User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d machines", len(*machines)),
		data.Channel, conn.UserID, "machine", 0, false)

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

	DB.Select(`id,name,photo_thumbnail,
		electrical,hydraulic,pnuematic,lube,printer,console,uncoiler,rollbed,conveyor`).
		From(`machine_type`).OrderBy(`name`).QueryStructs(machineTypes)

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

	DB.SQL(`select count(*) as num_tools from machine_type_tool where machine_id=$1`, data.ID).
		QueryScalar(&machineType.NumTools)

	logger(start, "Machine.GetMachineType",
		fmt.Sprintf("Channel %d, ID %d User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		machineType.Name,
		data.Channel, conn.UserID, "machine_type", data.ID, true)

	return nil
}

func (m *MachineRPC) UpdateMachineType(data shared.MachineTypeRPCData, done *bool) error {
	start := time.Now()

	// log.Println("here", data.MachineType)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.Update("machine_type").
		SetWhitelist(data.MachineType,
			"name",
			"electrical", "hydraulic", "pnuematic",
			"console", "printer", "lube",
			"uncoiler", "rollbed", "conveyor").
		Where("id = $1", data.ID).
		Exec()

	logger(start, "Machine.UpdateMachineType",
		fmt.Sprintf("Channel %d, ID %d User %d %s %s",
			data.Channel, data.ID, conn.UserID, conn.Username, conn.UserRole),
		data.MachineType.Name,
		data.Channel, conn.UserID, "machine_type", data.ID, true)

	*done = true
	return nil
}

func (m *MachineRPC) DeleteMachineType(data shared.MachineTypeRPCData, done *bool) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.SQL(`delete from machine_type where id=$1`, data.ID).Exec()

	logger(start, "Machine.DeleteMachineType",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("ID: %d", data.ID),
		data.Channel, conn.UserID, "machine_type", data.ID, true)

	*done = true
	return nil
}

func (m *MachineRPC) InsertMachineType(data shared.MachineTypeRPCData, id *int) error {
	start := time.Now()

	// log.Println("here", data)
	conn := Connections.Get(data.Channel)
	// log.Println("conn", conn)

	DB.InsertInto("machine_type").
		Columns("name",
			"electrical", "hydraulic", "pnuematic",
			"console", "printer", "lube",
			"uncoiler", "rollbed", "conveyor").
		Record(data.MachineType).
		Returning("id").
		QueryScalar(id)

	logger(start, "Machine.UpdateMachineType",
		fmt.Sprintf("Channel %d, ID %d User %d %s %s",
			data.Channel, *id, conn.UserID, conn.Username, conn.UserRole),
		data.MachineType.Name,
		data.Channel, conn.UserID, "machine_type", *id, true)

	return nil
}

func (m *MachineRPC) MachineTypeTools(data shared.MachineTypeRPCData, machineTypeTools *[]shared.MachineTypeTool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Select(`*`).
		From(`machine_type_tool`).
		OrderBy(`position`).
		Where(`machine_id=$1`, data.ID).
		QueryStructs(machineTypeTools)

	logger(start, "Machine.MachineTypes",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d machine type tools", len(*machineTypeTools)),
		data.Channel, conn.UserID, "machine_type", 0, true)

	return nil
}

func (m *MachineRPC) GetMachineTypeTool(data shared.MachineTypeToolRPCData, machineTypeTool *shared.MachineTypeTool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Select(`*`).
		From(`machine_type_tool`).
		OrderBy(`position`).
		Where(`machine_id=$1 and position=$2`, data.MachineID, data.ID).
		QueryStruct(machineTypeTool)

	logger(start, "Machine.GetMachineTypeTool",
		fmt.Sprintf("Channel %d, Machine %d Tool %d User %d %s %s",
			data.Channel, data.MachineID, data.ID, conn.UserID, conn.Username, conn.UserRole),
		machineTypeTool.Name,
		data.Channel, conn.UserID, "machine_type_tool", data.MachineID, false)

	return nil
}

func (m *MachineRPC) DeleteMachineTypeTool(data shared.MachineTypeToolRPCData, done *bool) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	*done = false

	DB.SQL(`delete from machine_type_tool where machine_id=$1 and position=$2`,
		data.MachineID, data.ID).Exec()

	logger(start, "Machine.DeleteMachineTypeTool",
		fmt.Sprintf("Channel %d, Machine %d Tool %d User %d %s %s",
			data.Channel, data.MachineID, data.ID, conn.UserID, conn.Username, conn.UserRole),
		"Deleted",
		data.Channel, conn.UserID, "machine_type_tool", data.MachineID, true)

	*done = true

	return nil
}

func (m *MachineRPC) UpdateMachineTypeTool(data shared.MachineTypeToolRPCData, done *bool) error {
	start := time.Now()

	// log.Println("here", data.MachineType)
	conn := Connections.Get(data.Channel)

	DB.Update("machine_type_tool").
		SetWhitelist(data.MachineTypeTool, "name").
		Where("machine_id = $1 and position=$2", data.MachineID, data.ID).
		Exec()

	logger(start, "Machine.UpdateMachineTypeTool",
		fmt.Sprintf("Channel %d, ID %d %d User %d %s %s",
			data.Channel, data.MachineID, data.ID, conn.UserID, conn.Username, conn.UserRole),
		data.MachineTypeTool.Name,
		data.Channel, conn.UserID, "machine_type_tool", data.ID, true)

	*done = true
	return nil
}
