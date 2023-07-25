package util

import "github.com/google/uuid"

func GetId() string {
	return uuid.NewString()
}
