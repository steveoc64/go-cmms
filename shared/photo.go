package shared

type Photo struct {
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
