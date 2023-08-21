package dto

import "github.com/shopspring/decimal"

type UserDto struct {
	Fullname string `json:"name"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

type ServiceDto struct {
	Type        uint8           `json:"type"`
	State       uint8           `json:"state"`
	Currency    string          `json:"currency"`
	InitBalance decimal.Decimal `json:"init_balance"`
}
