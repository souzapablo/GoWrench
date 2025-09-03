package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/connection_settings"
	"wrench/app/manifest/validation"
)

func dynamodbCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if appSetting.Connections != nil &&
		appSetting.Connections.DynamoDb != nil {

		result.Append(dynamodbTableIdDuplicated(appSetting.Connections.DynamoDb))
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
