package photo

import "github.com/Calmantara/go-user/common/entity"

type Photo struct {
	ID     uint64 `json:"-"`
	UserId uint64 `json:"-"`
	Name   string `json:"name"`
	entity.DefaultColumn
}
