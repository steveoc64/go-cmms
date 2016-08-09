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
	DB.SQL(`select
		* from component c
		left join machine_type_tool x on x.id=c.mtt_id
		where c.machine_id=$1
		order by x.position,c.zindex,lower(c.name)`, data.ID).
		QueryStructs(&machine.Components)

	// err = DB.Select("*").
	// 	From("component").
	// 	Where("machine_id = $1", data.ID).
	// 	OrderBy("position,zindex,lower(name)").
	// 	QueryStructs(&machine.Components)

	// fetch some basic info, flags and thumbnail from the parent machine type
	DB.Select(`name,photo_thumbnail,electrical,hydraulic,pnuematic,lube,printer,console,uncoiler,rollbed,conveyor,encoder,strip_guide`).
		From(`machine_type`).
		Where(`id=$1`, machine.MachineType).
		QueryStruct(&machine.MachineTypeData)

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
		// fetch all components
		DB.SQL(`select
		* from component c
		left join machine_type_tool x on x.id=c.mtt_id
		where c.machine_id=$1
		order by x.position,c.zindex,lower(c.name)`, m.ID).
			QueryStructs(&(*machines)[k].Components)

		// err = DB.Select("*").
		// 	From("component").
		// 	Where("machine_id = $1", m.ID).
		// 	OrderBy("position,zindex,lower(name)").
		// 	QueryStructs(&(*machines)[k].Components)
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
			"alerts_to", "tasks_to", "machine_type").
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
			"alerts_to", "tasks_to", "machine_type").
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
		electrical,hydraulic,pnuematic,lube,printer,console,uncoiler,rollbed,conveyor,encoder,strip_guide`).
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

	DB.Select(`id,name,photo_preview,
		electrical,hydraulic,pnuematic,lube,printer,console,uncoiler,rollbed,conveyor,encoder,strip_guide`).
		From(`machine_type`).
		Where(`id=$1`, data.ID).
		QueryStruct(machineType)

	// fetch the tool count
	DB.SQL(`select count(*) as num_tools from machine_type_tool where machine_id=$1`, data.ID).
		QueryScalar(&machineType.NumTools)

	// fetch the actual tools
	DB.Select(`position,name`).
		From(`machine_type_tool`).
		Where(`machine_id=$1`, data.ID).
		OrderBy(`position`).
		QueryStructs(&machineType.Tools)

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
			"uncoiler", "rollbed", "conveyor",
			"encoder", "strip_guide").
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
			"uncoiler", "rollbed", "conveyor",
			"encoder", "strip_guide").
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
		Where(`id=$1`, data.ID).
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

	oldPos := 0
	DB.SQL(`select position from machine_type_tool where id=$1`, data.ID).
		QueryScalar(&oldPos)

	DB.SQL(`delete from machine_type_tool where id=$1`, data.ID).Exec()

	// and now shuffle everything down by one from this point
	DB.SQL(`update machine_type_tool set position=(position-1) where machine_id=$1 and position > $2`,
		data.MachineID, oldPos).
		Exec()

	logger(start, "Machine.DeleteMachineTypeTool",
		fmt.Sprintf("Channel %d, Machine %d Tool %d User %d %s %s",
			data.Channel, data.MachineID, data.ID, conn.UserID, conn.Username, conn.UserRole),
		"Deleted",
		data.Channel, conn.UserID, "machine_type_tool", data.MachineID, true)

	*done = true

	rehashTools(data.MachineID)

	return nil
}

func (m *MachineRPC) UpdateMachineTypeTool(data shared.MachineTypeToolRPCData, done *bool) error {
	start := time.Now()

	// log.Println("here", data.MachineType)
	conn := Connections.Get(data.Channel)

	// Get the original position
	oldPos := 0
	DB.SQL(`select position from machine_type_tool where id=$1`, data.ID).QueryScalar(&oldPos)
	println("Pos:", oldPos, data.MachineTypeTool.Position)

	if oldPos != data.MachineTypeTool.Position {
		println("need to reshuffle position from", oldPos, "to", data.MachineTypeTool.Position)

		// shuffle everything down from the OLD position
		DB.SQL(`update machine_type_tool set position=(position-1) where machine_id=$1 and position > $2 and id != $3`,
			data.MachineID, oldPos, data.ID).Exec()

		// insert a new space for this NEW position
		DB.SQL(`update machine_type_tool set position=(position+1) where machine_id=$1 and position >= $2 and id != $3`,
			data.MachineID, data.MachineTypeTool.Position, data.ID).Exec()
	}

	DB.Update("machine_type_tool").
		SetWhitelist(data.MachineTypeTool, "name", "position").
		Where("id=$1", data.ID).
		Exec()

	logger(start, "Machine.UpdateMachineTypeTool",
		fmt.Sprintf("Channel %d, ID %d %d User %d %s %s",
			data.Channel, data.MachineID, data.ID, conn.UserID, conn.Username, conn.UserRole),
		data.MachineTypeTool.Name,
		data.Channel, conn.UserID, "machine_type_tool", data.ID, true)

	rehashTools(data.MachineID)

	*done = true
	return nil
}

func (m *MachineRPC) InsertMachineTypeTool(data shared.MachineTypeToolRPCData, id *int) error {
	start := time.Now()
	*id = 0

	// log.Println("here", data.MachineType)
	conn := Connections.Get(data.Channel)

	// If there is already a record at this position, then shuffle them all up from here on
	DB.SQL(`update machine_type_tool set position=(position+1) where machine_id=$1 and position >= $2`,
		data.MachineTypeTool.MachineID,
		data.MachineTypeTool.Position).
		Exec()

	DB.InsertInto("machine_type_tool").
		Columns("machine_id", "position", "name").
		Record(data.MachineTypeTool).
		Exec()
	*id = data.MachineTypeTool.ID

	logger(start, "Machine.InsertMachineTypeTool",
		fmt.Sprintf("Channel %d, User %d %s %s",
			data.Channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d %v", *id, data.MachineTypeTool),
		data.Channel, conn.UserID, "machine_type_tool", *id, true)

	rehashTools(data.MachineTypeTool.MachineID, *id, "insert", &data)

	return nil
}

func rehashTools(mt int, mtt int, mode string, data *shared.MachineTypeTool) {
	println("rehashTools", mt, mtt, mode)

	switch mode {
	case "insert":
		// Need to create a whole new component record for each machine instance of the same machinetype
		comp := shared.Component{
			MachineID: data.MachineID,
			SiteID: data.S
			machine_id | integer | not null
 id         | integer | not null default nextval('component_id_seq'::regclass)
 site_id    | integer | not null
 name       | text    | not null
 descr      | text    | not null default ''::text
 make       | text    | not null default ''::text
 model      | text    | not null default ''::text
 serialnum  | text    | not null default ''::text
 notes      | text    | not null default ''::text
 qty        | integer | not null default 1
 stock_code | text    | not null default ''::text
 position   | integer | not null default 1
 status     | text    | not null default 'Running'::text
 is_running | boolean | not null default true
 zindex     | integer | not null default 0
 mtt_id     | integer | not null default 0

		}
	case "delete":
	case "update":
	}
}
