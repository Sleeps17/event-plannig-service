package models

type Employee struct {
	ID           uint64 `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Patronymic   string `json:"patronymic"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
}
