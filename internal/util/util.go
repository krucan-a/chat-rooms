package util

import (
	"fmt"
	"github.com/segmentio/ksuid"
)

func GenerateID() string {
	return ksuid.New().String()
}

func GenerateCookieName(roomID string) string {
	return fmt.Sprintf("%s+%s", "user-id", roomID)
}
