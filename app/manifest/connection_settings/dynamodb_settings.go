package connection_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type DynamoDbConnectionSettings struct {
	Local                   bool   `yaml:"local"`
	LocalEndpoint           string `yaml:"localEndpoint"`
	LocalAwsAccessKeyId     string `yaml:"localAwsAccessKeyId"`
	LocalAwsSecretAccessKey string `yaml:"localAwsSecretAccessKey"`
	LocalAwsRegion          string `yaml:"localAwsRegion"`

	Tables []*DynamoDbTableSettings `yaml:"tables"`
}

func (setting *DynamoDbConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Tables) == 0 {
		result.AddError("connections.dynamodb.tables is required")
	} else {
		for _, s := range setting.Tables {
			result.AppendValidable(s)
		}
	}

	if setting.Local {

		if len(setting.LocalEndpoint) == 0 {
			result.AddError("connections.dynamodb.localEndpoint is required when local is true")
		}
		if len(setting.LocalAwsAccessKeyId) == 0 {
			result.AddError("connections.dynamodb.localAwsAccessKeyId is required when local is true")
		}
		if len(setting.LocalAwsSecretAccessKey) == 0 {
			result.AddError("connections.dynamodb.localAwsSecretAccessKey is required when local is true")
		}
		if len(setting.LocalAwsRegion) == 0 {
			result.AddError("connections.dynamodb.localAwsRegion is required when local is true")
		}
	}

	return result
}

type DynamoDbTableSettings struct {
	Id               string `yaml:"id"`
	Name             string `yaml:"name"`
	PartitionKeyName string `yaml:"partitionKeyName"`
	SortKeyName      string `yaml:"sortKeyName"`
}

func (setting *DynamoDbTableSettings) GetId() string {
	return setting.Id
}

func (setting *DynamoDbTableSettings) Valid() validation.ValidateResult {
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
