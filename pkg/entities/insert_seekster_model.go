package entities

type InsertSeeksterUserInput struct {
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(15)" validate:"required,numeric,gte=10,lte=15"`
	Password    string `json:"password" validate:"required"`
	UUID        string `json:"uuid" validate:"required"`
	SSOID       string `json:"ssoid" validate:"required"`
}
