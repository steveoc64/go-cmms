package shared

type RouteReq struct {
	Channel int
	Name    string
}

type RouteResponse struct {
	Template string
}
