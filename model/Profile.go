package model

type Profile struct {
	Id      uint   `json:"id" gorm:"primary_key"`
	UserId  uint   `json:"user_id"`
	Profile string `json:"profile"`
}
