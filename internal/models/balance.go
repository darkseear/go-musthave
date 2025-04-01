package models

type Balance struct {
	Current   float64 `json:"balance" db:"balance"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}
