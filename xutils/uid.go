package xutils

import (
	"encoding/json"
	"github.com/google/uuid"
	//uuid "github.com/satori/go.uuid"
)

// GenerateUUID generates a new UUID string.
func GenerateUUID() string {
	// Generate a new UUID using the google/uuid package
	uuidWithHyphen := uuid.New()
	return uuidWithHyphen.String()
	//return uuid.NewV4().String()
}

func JsonStr(v any) string {
	d, e := json.Marshal(v)
	if e != nil {
		return ""
	} else {
		return string(d)
	}
}
