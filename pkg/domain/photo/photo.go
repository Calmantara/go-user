package photo

import (
	"github.com/Calmantara/go-user/lib/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Photo struct {
	ID     uint64 `json:"-"`
	UserId uint64 `json:"-"`
	Name   string `json:"name"`
	entity.DefaultColumn
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "name"},
				{Name: "user_id"}},
			DoNothing: true},
	)
	return tx.Error
}
