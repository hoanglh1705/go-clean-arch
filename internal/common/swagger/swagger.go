package swagger

type Ok struct{}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"trace_id"`
}
