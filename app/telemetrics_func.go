package app

import (
	"strings"

	"github.com/google/uuid"
)

var instanceID string

func LoadInstanceID() {
	instance := uuid.New()
	instanceId := instance.String()
	instanceId = strings.ReplaceAll(instanceId, "-", "")

	instanceID = instanceId
}

func GetInstanceID() string {
	return instanceID
}
