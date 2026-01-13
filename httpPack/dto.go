package httppack

import (
	"encoding/json"
	"time"
)

type WorkersDTO struct {
	Id       int
	FullName string
	Position string
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
