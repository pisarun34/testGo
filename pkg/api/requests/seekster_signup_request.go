package requests

type SignUpInput struct {
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(15)" validate:"required,numeric,gte=10,lte=15"`
	Ssoid       string `json:"ssoid" gorm:"type:varchar(15)" validate:"required,numeric,gte=10,lte=15"`
}
