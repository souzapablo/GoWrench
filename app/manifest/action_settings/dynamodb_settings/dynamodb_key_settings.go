package dynamodb_settings

import "wrench/app/manifest/validation"

type DynamoDbKeySettings struct {
	PartitionKeyValue string `yaml:"partitionKeyValue"`
	SortKeyValue      string `yaml:"sortKeyValue"`
}

func (settings *DynamoDbKeySettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.PartitionKeyValue) == 0 {
		result.AddError("actions.dynamodb.key.partitionKeyValue is required")
	}

	return result
}
