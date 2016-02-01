package shared

type PingReq struct {
	PingName string
	Seq      int
}

type PingRep struct {
	Seq    int
	Result string
}
