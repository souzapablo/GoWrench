package key_settings

import "wrench/app/manifest/validation"

type KeySettings struct {
	Id                     string `yaml:"id"`
	PrivateRsaKeyDERBase64 string `yaml:"privateRsaKeyDERBase64"`
	Passphrase             string `yaml:"passphrase"`
}

func (setting *KeySettings) GetId() string {
	return setting.Id
}

func (setting KeySettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult
	if len(setting.Id) == 0 {
		result.AddError("keySettings.id is required")
	}
	if len(setting.PrivateRsaKeyDERBase64) == 0 {
		result.AddError("keySettings.privateRsaKeyDERBase64 is required")
	}

	return result
}
