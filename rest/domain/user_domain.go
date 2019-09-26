package domain

type User struct {
	Id int64 `json:"id"`
	UserId string `json:"userId"`
	Password string `json:"password"`
	Name string `json:"name"`
	Gender string `json:"genderr"`
	PhoneNumber string `json:"phoneNumber"`
	Country string `json:"country"`
	Language string `json:"language"`
	Remark string `json:"remark"`
}