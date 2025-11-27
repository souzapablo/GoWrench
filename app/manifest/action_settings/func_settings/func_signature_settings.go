package func_settings

import (
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type FuncSignatureSettings struct {
	KeyId     string        `yaml:"keyId"`
	Algorithm types.HashAlg `yaml:"algorithm"`
}

func (setting FuncSignatureSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult
	if len(setting.KeyId) == 0 {
		result.AddError("actions.func.sign.keyId is required")
	}

	if len(setting.Algorithm) == 0 {
		result.AddError("actions.func.sign.algorithm is required")
	}

	return result
}
