package dynamodb_settings

import "wrench/app/manifest/validation"

type DynamoDbCommand string

const (
	DynamoDbCommandCreate         DynamoDbCommand = "create"
	DynamoDbCommandUpdate         DynamoDbCommand = "update"
	DynamoDbCommandCreateOrUpdate DynamoDbCommand = "createOrUpdate"
	DynamoDbCommandDelete         DynamoDbCommand = "delete"
	DynamoDbCommandGet            DynamoDbCommand = "get"
	DynamoDbCommandList           DynamoDbCommand = "list"
)

type DynamoDbSettings struct {
	TableId string               `yaml:"tableId"`
	Command DynamoDbCommand      `yaml:"command"`
	Key     *DynamoDbKeySettings `yaml:"key"`
}

func (settings *DynamoDbSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.TableId) == 0 {
		result.AddError("actions.dynamodb.tableId is required")
	}

	if len(settings.Command) == 0 {
		result.AddError("actions.dynamodb.command is required")
	} else {
		if (settings.Command == DynamoDbCommandCreate ||
			settings.Command == DynamoDbCommandUpdate ||
			settings.Command == DynamoDbCommandCreateOrUpdate ||
			settings.Command == DynamoDbCommandDelete ||
			settings.Command == DynamoDbCommandGet ||
			settings.Command == DynamoDbCommandList) == false {
			result.AddError("ctions.dynamodb.command should contain valid value (create, update, createOrUpdate, delete, get or list)")
		}

		if settings.Command == DynamoDbCommandGet ||
			settings.Command == DynamoDbCommandDelete {

			if settings.Key != nil {
				result.Append(settings.Key.Valid())
			} else {
				result.AddError("actions.dynamodb.key is required when command is get or delete")
			}
		}
	}

	return result
}
