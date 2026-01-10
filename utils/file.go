package utils

import (
	"os"
)

func GetSuperuserToken() []byte {
	token, err := os.ReadFile(".superusertoken")
	if err != nil {
		panic("Cannot Read SuperuserToken File" + err.Error())
	}
	return token
}
