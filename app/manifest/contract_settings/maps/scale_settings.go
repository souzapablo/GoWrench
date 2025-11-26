package maps

import (
	"strings"
	"wrench/app/manifest/validation"
)

type ScaleSettings struct {
	Up   []string `yaml:"up"`
	Down []string `yaml:"down"`
}

func (setting ScaleSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Up) == 0 && len(setting.Down) == 0 {
		result.AddError("contract.maps.scale must configure up or down")
		return result
	}

	validateList := func(list []string, op string) {
		for _, property := range list {
			errMsg := "contract.maps.scale." + op + " should be configured as 'property:number' without spaces"

			if strings.Contains(property, " ") {
				result.AddError(errMsg)
				continue
			}

			parts := strings.Split(property, ":")
			if len(parts) != 2 {
				result.AddError(errMsg)
				continue
			}
		}
	}

	validateList(setting.Up, "up")
	validateList(setting.Down, "down")

	return result
}
