package shared

import "github.com/steveoc64/formulate"

type Photo struct {
	ID       int    `db:"id"`
	Type     string `db:"type"`
	Datatype string `db:"datatype"`
	Filename string `db:"filename"`
	Entity   string `db:"entity"`
	EntityID int    `db:"entity_id"`
	Data     string `db:"photo"`
	Preview  string `db:"preview"`
	Thumb    string `db:"thumb"`
	Notes    string `db:"notes"`
}

type Phototest struct {
	ID        int                 `db:"id"`
	Notes     string              `db:"name"`
	Photo     formulate.FileField `db:"photo"`
	Preview   string              `db:"preview"`
	Thumbnail string              `db:"thumbnail"`
}

type PhotoTestRPCData struct {
	Channel int
	ID      int
	Photo   *Phototest
}

type PhotoRPCData struct {
	Channel int
	ID      int
	Photo   *Photo
}
