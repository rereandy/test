package model

import "time"

type Person struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Age   uint64 `json:"age"`
}

type CreatePersonRequest struct {
	Person Person `json:"person"`
}

type CreatePersonResponse struct {
	Person   Person        `json:"person"`
	Elapse   time.Duration `json:"elapse"` // nano seconds
	BaseResp BaseResp      `json:"baseresp"`
}
