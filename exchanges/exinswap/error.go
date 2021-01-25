package exinswap

import (
	"fmt"
)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"message,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}
