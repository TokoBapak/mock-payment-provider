package schema

type Error struct {
	StatusCode    int    `json:"status_code,string"`
	StatusMessage string `json:"status_message"`
}
