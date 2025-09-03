package connection_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type DynamodbConnectionSettings struct {
	Local  bool                     `yaml:"local"`
	Tables []*DynamodbTableSettings `yaml:"tables"`
}

func (setting *DynamodbConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Tables) == 0 {
		result.AddError("connections.dynamodb.tables is required")
	} else {
		for _, s := range setting.Tables {
			result.AppendValidable(s)
		}
	}

	return result
}

type DynamodbTableSettings struct {
	Id               string `yaml:"id"`
	Name             string `yaml:"name"`
	PartitionKeyName string `yaml:"partitionKeyName"`
	SortKeyName      string `yaml:"sortKeyName"`
}

func (setting *DynamodbTableSettings) GetId() string {
	return setting.Id
}

func (setting *DynamodbTableSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("connections.dynamodb.id is required")
	}

	if len(setting.Name) == 0 {
		result.AddError(fmt.Sprintf("connections.dynamodb.tables[%v].name is required", setting.Id))
	}

	if len(setting.PartitionKeyName) == 0 {
		result.AddError(fmt.Sprintf("connections.dynamodb.tables[%v].partitionKeyName is required", setting.Id))
	}

	return result
}
