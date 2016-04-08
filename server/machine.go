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

	logger(start, "Machine.Get",
		fmt.Sprintf("Machine %d", machineID),
		machine.Name)

	return nil
}
