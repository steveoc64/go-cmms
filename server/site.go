package main

import (
	"fmt"
	"log"
	"time"

	"github.com/steveoc64/go-cmms/shared"
)

type SiteRPC struct{}

const SiteQueryAll = `select 
s.*,p.name as parent_site_name,t.name as stock_site_name
		from site s
		left join site p on (p.id=s.parent_site)
		left join site t on (t.id=s.stock_site)
order by lower(s.name)`

const SiteQueryBySite = `select 
s.*,p.name as parent_site_name,t.name as stock_site_name
		from site s
		left join site p on (p.id=s.parent_site)
		left join site t on (t.id=s.stock_site)
where s.id = $1
order by lower(s.name)`

const SiteQueryInSite = `select 
s.*,p.name as parent_site_name,t.name as stock_site_name
		from site s
		left join site p on (p.id=s.parent_site)
		left join site t on (t.id=s.stock_site)
where s.id in $1
order by lower(s.name)`

const MachinesBySite = `select 
m.*,s.name as site_name,x.span as span
		from machine m
		left join site s on (s.id=m.site_id)
		left join site_layout x on (x.site_id=m.site_id and x.machine_id=m.id)
where m.site_id = $1
order by x.seq,lower(m.name)`

// How many sites does this user have ?
func (s *SiteRPC) SiteCount(channel int, count *int) error {
	start := time.Now()

	conn := Connections.Get(channel)

	*count = 0
	err := DB.SQL(`select count(*) from user_site where user_id=$1`, conn.UserID).QueryScalar(count)
	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Site.SiteCount",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Sites", count))

	return nil
}

// Get a list of all sites, filtered by User
func (s *SiteRPC) UserList(channel int, sites *[]shared.Site) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Read the DB to get a list of sites this user has access to
	userSites := []int{}
	DB.SQL(`select site_id from user_site where user_id=$1`, conn.UserID).QuerySlice(&userSites)

	// Read the sites that this user has access to
	err := DB.SQL(SiteQueryInSite, userSites).QueryStructs(sites)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Site.UserList",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Sites", len(*sites)))

	return nil
}

// Get a list of all sites, which is mainly used for lookup purposes
func (s *SiteRPC) List(channel int, sites *[]shared.Site) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Read all the sites
	err := DB.SQL(SiteQueryAll).QueryStructs(sites)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Site.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Sites", len(*sites)))

	return nil
}

// Get the details for a given site
func (s *SiteRPC) Get(siteID int, site *shared.Site) error {
	start := time.Now()

	// Read the sites that this user has access to
	err := DB.SQL(SiteQueryBySite, siteID).QueryStruct(site)

	if err != nil {
		log.Println(err.Error())
	}

	// bar := "==============================================\n"
	logger(start, "Site.Get",
		fmt.Sprintf("Site %d", siteID),
		site.Name)
	// fmt.Sprintf("%s\n%s%#v\n%s", site.Name, bar, site, bar))

	return nil
}

// Save a site
func (s *SiteRPC) Save(data shared.SiteUpdateData, retval *int) error {
	start := time.Now()

	conn := Connections.Get(data.Channel)

	DB.Update("site").
		SetWhitelist(data.Site, "name", "address", "phone", "fax", "parent_site", "stock_site", "notes").
		Where("id = $1", data.Site.ID).
		Exec()

	logger(start, "Site.Save",
		fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
			data.Channel, data.Site.ID, conn.UserID, conn.Username, conn.UserRole),
		data.Site.Name)

	return nil
}

// Get the details for my home site
func (s *SiteRPC) GetHome(channel int, site *shared.Site) error {
	start := time.Now()

	// Get my home site id
	conn := Connections.Get(channel)
	siteID := 0
	DB.SQL(`select site_id from user_site where user_id=$1 limit 1`, conn.UserID).QueryScalar(&siteID)

	// Read the sites that this user has access to
	err := DB.SQL(SiteQueryBySite, siteID).QueryStruct(site)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "Site.GetHome",
		fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
			channel, siteID, conn.UserID, conn.Username, conn.UserRole),
		site.Name)

	return nil
}

// Get all machines for the given site
func (s *SiteRPC) MachineList(req *shared.MachineReq, machines *[]shared.Machine) error {
	start := time.Now()

	conn := Connections.Get(req.Channel)

	// Read the machines for the given site
	err := DB.SQL(MachinesBySite, req.SiteID).QueryStructs(machines)

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

	logger(start, "Site.MachineList",
		fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
			req.Channel, req.SiteID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d machines", len(*machines)))

	return nil
}

// Get all machines for the home site
func (s *SiteRPC) HomeMachineList(channel int, machines *[]shared.Machine) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Get my home site
	siteID := 0
	DB.SQL(`select site_id from user_site where user_id=$1 limit 1`, conn.UserID).QueryScalar(&siteID)

	// Read the machines for the given site
	err := DB.SQL(MachinesBySite, siteID).QueryStructs(machines)

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

	logger(start, "Site.HomeMachineList",
		fmt.Sprintf("Channel %d, Site %d, User %d %s %s",
			channel, siteID, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d machines", len(*machines)))

	return nil
}

// Get a SiteStatus Report
func (s *SiteRPC) StatusReport(channel int, retval *shared.SiteStatusReport) error {
	start := time.Now()

	conn := Connections.Get(channel)

	retval.Edinburgh = "Running"
	retval.Minto = "Running"
	retval.Tomago = "Running"
	retval.Chinderah = "Running"

	i := 0

	// Get the overall status for Edinburgh
	DB.SQL(`select count(m.*) 
		from machine m
		left join site s on (s.id = m.site_id)
		where m.status = 'Stopped' 
		and s.name like 'Edinburgh%'`).QueryScalar(&i)
	if i > 0 {
		retval.Edinburgh = "Stopped"
	} else {
		DB.SQL(`select count(m.*) 
			from machine m
			left join site s on (s.id = m.site_id)
			where m.status = 'Needs Attention' 
			and s.name like 'Edinburgh%'`).QueryScalar(&i)
		if i > 0 {
			retval.Edinburgh = "Needs Attention"
		}
	}

	// Get the overall status for Minto
	i = 0
	DB.SQL(`select count(m.*) 
		from machine m
		left join site s on (s.id = m.site_id)
		where m.status = 'Stopped' 
		and s.name = 'Minto'`).QueryScalar(&i)
	if i > 0 {
		retval.Minto = "Stopped"
	} else {
		DB.SQL(`select count(m.*) 
			from machine m
			left join site s on (s.id = m.site_id)
			where m.status = 'Needs Attention' 
			and s.name = 'Minto'`).QueryScalar(&i)
		if i > 0 {
			retval.Minto = "Needs Attention"
		}
	}

	// Get the overall status for Tomago
	i = 0
	DB.SQL(`select count(m.*) 
		from machine m
		left join site s on (s.id = m.site_id)
		where m.status = 'Stopped' 
		and s.name = 'Tomago'`).QueryScalar(&i)
	if i > 0 {
		retval.Tomago = "Stopped"
	} else {
		DB.SQL(`select count(m.*) 
			from machine m
			left join site s on (s.id = m.site_id)
			where m.status = 'Needs Attention' 
			and s.name = 'Tomago'`).QueryScalar(&i)
		if i > 0 {
			retval.Tomago = "Needs Attention"
		}
	}

	// Get the overall status for Chinderah
	i = 0
	DB.SQL(`select count(m.*) 
		from machine m
		left join site s on (s.id = m.site_id)
		where m.status = 'Stopped' 
		and s.name = 'Chinderah'`).QueryScalar(&i)
	if i > 0 {
		retval.Chinderah = "Stopped"
	} else {
		DB.SQL(`select count(m.*) 
			from machine m
			left join site s on (s.id = m.site_id)
			where m.status = 'Needs Attention' 
			and s.name = 'Chinderah'`).QueryScalar(&i)
		if i > 0 {
			retval.Chinderah = "Needs Attention"
		}
	}

	logger(start, "Site.StatusReport",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("E: %s M: %s T: %s C: %s",
			retval.Edinburgh,
			retval.Minto,
			retval.Tomago,
			retval.Chinderah))

	return nil
}