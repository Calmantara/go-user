package photo

type Photo struct {
	ID     uint64 `json:"id"`
	UserId uint64 `json:"user_id"`
	Name   string `json:"name"`
}
