package cross_funcs

import (
	"strings"

	"github.com/google/uuid"
)

var instanceID string

func init() {
	instanceID = loadInstanceID()
}

func loadInstanceID() string {
	instance := uuid.New()
	instanceId := instance.String()
	instanceId = strings.ReplaceAll(instanceId, "-", "")

	return instanceId
}

func GetInstanceID() string {
	return instanceID
}
