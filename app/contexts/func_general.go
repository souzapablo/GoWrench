package contexts

import (
	"encoding/base64"
	"strconv"
	"time"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/func_settings"
)

func GetFuncValue(funcType func_settings.FuncGeneralType, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) (string, error) {

	body, err := bodyContext.GetBody(action)
	if err != nil {
		return "", err
	}

	switch funcType {
	case func_settings.FuncTypeTimestampMilli:
		return getTimestamp(), nil
	case func_settings.FuncTypeBase64Encode:
		return base64.StdEncoding.EncodeToString(body), nil
	case func_settings.FuncTypeBase64UrlEncode:
		return base64.RawURLEncoding.EncodeToString(body), nil
	default:
		return "", nil
	}
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
}
