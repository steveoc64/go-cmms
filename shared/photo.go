package shared

type Photo struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Photo     string `db:"photo"`
	Thumbnail string `db:"thumbnail"`
	Photodec  string `db:"photodec"`
}

type PhotoRPCData struct {
	Channel int
	ID      int
	Photo   *Photo
}
