package main

import (
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type SiteRPC struct{}

func (s *SiteRPC) List(channel int, sites *[]shared.Site) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Read the DB to get a list of sites this user has access to
	userSites := []int{}
	DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&userSites)

	// Read the sites that this user has access to
	err := DB.SQL(`select s.*,p.name as parent_site_name,t.name as stock_site_name
		from site s
		left join site p on (p.id=s.parent_site)
		left join site t on (t.id=s.stock_site)
		where s.id in $1
		order by lower(s.name)`, userSites).QueryStructs(sites)

	if err != nil {
		log.Println(err.Error())
	}

	log.Printf(`Site.List -> %s
    » (Channel %d, User %d %s %s)
    « (%d sites)`,
		time.Since(start),
		channel, conn.UserID, conn.Username, conn.UserRole,
		len(*sites))

	return nil
}

func (s *SiteRPC) Get(siteID int, site *shared.Site) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(`select s.*,p.name as parent_site_name,t.name as stock_site_name
		from site s
		left join site p on (p.id=s.parent_site)
		left join site t on (t.id=s.stock_site)
		where s.id = $1
		order by lower(s.name)`, siteID).QueryStruct(site)

	if err != nil {
		log.Println(err.Error())
	}

	log.Printf(`Site.Get -> %s
    » (Site %d)
    « (%s ...)`,
		time.Since(start),
		siteID,
		site.Name)

	return nil
}

func (s *SiteRPC) MachineList(req *shared.MachineReq, machines *[]shared.Machine) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	// Read the machines for the given site
	err := DB.SQL(`select m.*,s.name as site_name,x.span as span
		from machine m
		left join site s on (s.id=m.site_id)
		left join site_layout x on (x.site_id=m.site_id and x.machine_id=m.id)
		where m.site_id = $1
		order by x.seq,lower(m.name)`, req.SiteID).QueryStructs(machines)

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

	log.Printf(`Site.MachineList -> %s
    » (Channel %d, Site %d, User %d %s %s)
    « (%d machines)`,
		time.Since(start),
		req.Channel, req.SiteID, conn.UserID, conn.Username, conn.UserRole,
		len(*machines))

	return nil
}
