package dynamodb_settings

import "wrench/app/manifest/validation"

type DynamoDbCommand string

const (
	DynamoDbCommandCreate DynamoDbCommand = "create"
	DynamoDbCommandUpdate DynamoDbCommand = "update"
	DynamoDbCommandDelete DynamoDbCommand = "delete"
	DynamoDbCommandGet    DynamoDbCommand = "get"
	DynamoDbCommandList   DynamoDbCommand = "list"
)

type DynamodbSettings struct {
	TableId           string          `yaml:"tableId"`
	Command           DynamoDbCommand `yaml:"command"`
	PartitionKeyValue string          `yaml:"partitionKeyValue"`
	SortKeyValue      string          `yaml:"sortKeyValue"`
}

func (settings *DynamodbSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.TableId) == 0 {
		result.AddError("actions.dynamodb.tableId is required")
	}

	if len(settings.Command) == 0 {
		result.AddError("actions.dynamodb.command is required")
	} else {
		if (settings.Command == DynamoDbCommandCreate ||
			settings.Command == DynamoDbCommandUpdate ||
			settings.Command == DynamoDbCommandDelete ||
			settings.Command == DynamoDbCommandGet ||
			settings.Command == DynamoDbCommandList) == false {
			result.AddError("ctions.dynamodb.command should contain valid value (create, update, delete, get or list)")
		}
	}

	if len(settings.PartitionKeyValue) == 0 {
		result.AddError("actions.dynamodb.partitionKeyValue is required")
	}

	return result
}
