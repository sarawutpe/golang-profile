package model

type Avatar struct {
	Id     uint   `json:"id" gorm:"primary_key"`
	UserId uint   `json:"user_id"`
	Avatar string `json:"avatar"`
}
