package shared

type Photo struct {
	ID       int    `db:"id"`
	Type     string `db:"type"`
	Datatype string `db:"datatype"`
	Filename string `db:"filename"`
	Entity   string `db:"entity"`
	EntityID int    `db:"entity_id"`
	Photo    string `db:"photo"`
	Preview  string `db:"preview"`
	Thumb    string `db:"thumb"`
	Notes    string `db:"notes"`
}

type Phototest struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Photo     string `db:"photo"`
	Preview   string `db:"preview"`
	Thumbnail string `db:"thumbnail"`
}

type PhotoRPCData struct {
	Channel int
	ID      int
	Photo   *Photo
}
