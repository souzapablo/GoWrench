package cross_validation

import (
	"fmt"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/key_settings"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
)

func keyCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(appSetting.Keys) > 0 {
		result.Append(keyIdDuplicated(appSetting.Keys))
	}

	result.Append(keyIdRefExist(appSetting))

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

func keyIdRefExist(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	actions := getActionsByType(appSetting.Actions, action_settings.ActionTypeFuncSignature)

	if len(actions) > 0 {
		for _, action := range actions {
			if len(action.Func.Sign.KeyId) > 0 {
				_, err := manifest_cross_funcs.GetPrivateKeyById(action.Func.Sign.KeyId)

				if err != nil {
					result.AddError(fmt.Sprintf("actions[%s].func.sign.keyId.  Don't exist keyId %s informed", action.Id, action.Func.Sign.KeyId))
				}
			}
		}
	}

	return result
}
