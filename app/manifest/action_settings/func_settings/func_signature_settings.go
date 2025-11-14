package func_settings

import "wrench/app/manifest/validation"

type FuncSignatureSettings struct {
	KeyId string `yaml:"keyId"`
}

func (setting FuncSignatureSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult
	if len(setting.KeyId) == 0 {
		result.AddError("actions.func.sign.keyId is required")
	}
	return result
}
