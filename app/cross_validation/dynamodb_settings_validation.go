package cross_validation

import (
	"fmt"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/connection_settings"
	"wrench/app/manifest/validation"
	"wrench/app/startup/connections"
)

func dynamodbCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if appSetting.Connections != nil &&
		appSetting.Connections.DynamoDb != nil {

		result.Append(dynamodbTableIdDuplicated(appSetting.Connections.DynamoDb))
	}

	if len(appSetting.Actions) > 0 {
		result.Append(dynamoDbActionTableIdExist(appSetting.Actions))
	}

	return result
}

func dynamodbTableIdDuplicated(settings *connection_settings.DynamodbConnectionSettings) validation.ValidateResult {

	var result validation.ValidateResult

	if len(settings.Tables) > 0 {
		hasIds := toHasIdSlice(settings.Tables)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("connections.dynamodb.tables.id %v duplicated", id))
		}

	}

	return result
}

func dynamoDbActionTableIdExist(settings []*action_settings.ActionSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings) > 0 {
		for _, setting := range settings {
			if setting.DynamoDb != nil {
				_, err := connections.GetDynamoDbTableConnection(setting.DynamoDb.TableId)
				if err != nil {
					result.AddError(fmt.Sprintf("actions[%v].dynamodb.tableId don't exist in connections.dynamodb.tables", setting.Id))
				}
			}
		}
	}

	return result
}
