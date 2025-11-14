package key_settings

import "wrench/app/manifest/validation"

type KeySettings struct {
	Id                   string `yaml:"id"`
	PrivateRsaKeysBase64 string `yaml:"privateRsaKeysBase64"`
	Passphrase           string `yaml:"passphrase"`
}

func (setting *KeySettings) GetId() string {
	return setting.Id
}

func (setting KeySettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult
	if len(setting.Id) == 0 {
		result.AddError("keySettings.id is required")
	}
	if len(setting.PrivateRsaKeysBase64) == 0 {
		result.AddError("keySettings.privateRsaKeysBase64 is required")
	}

	return result
}
