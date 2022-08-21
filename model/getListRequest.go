package model

type GetListFeedback struct {
	ClientId    string
	DisplayName string
	CountryCode string `gorm:"varchar(40)"`
}
