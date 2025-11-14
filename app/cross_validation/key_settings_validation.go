package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/key_settings"
	"wrench/app/manifest/validation"
)

func keyCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(appSetting.Keys) > 0 {
		result.Append(keyIdDuplicated(appSetting.Keys))
	}

	return result
}

func keyIdDuplicated(settings []*key_settings.KeySettings) validation.ValidateResult {

	var result validation.ValidateResult

	hasIds := toHasIdSlice(settings)
	duplicateIds := duplicateIdsValid(hasIds)

	for _, id := range duplicateIds {
		result.AddError(fmt.Sprintf("keys.id %v duplicated", id))
	}

	return result
}
