package uuid

import (
	"strings"

	"github.com/google/uuid"
)

const (
	InstanceIdLength = 11
)

func GenerateWithLength(length int) string {
	uid := strings.ReplaceAll(uuid.NewString(), "-", "")
	if len(uid) > length {
		return uid[:length]
	}
	return uid
}
