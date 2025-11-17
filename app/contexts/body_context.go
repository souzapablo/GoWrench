package contexts

import (
	"encoding/json"
	"fmt"
	"strings"
	"wrench/app/manifest/action_settings"
)

type BodyContext struct {
	CurrentBodyByteArray []byte
	BodyPreserved        map[string][]byte
	HttpStatusCode       int
	ContentType          string
	Headers              map[string]string
}

func (bodyContext *BodyContext) SetBodyPreserved(id string, body []byte) {
	if bodyContext.BodyPreserved == nil {
		bodyContext.BodyPreserved = make(map[string][]byte)
	}

	bodyContext.BodyPreserved[id] = body
}

func (bodyContext *BodyContext) SetBody(body []byte) {
	bodyContext.CurrentBodyByteArray = body
}

func (bodyContext *BodyContext) SetBodyAction(settings *action_settings.ActionSettings, body []byte) {
	if settings.ShouldPreserveBody() {
		bodyContext.SetBodyPreserved(settings.Id, body)
	} else {
		bodyContext.SetBody(body)
	}
}

func (bodyContext *BodyContext) GetBodyPreserved(id string) ([]byte, error) {
	if bodyContext.BodyPreserved == nil {
		return nil, fmt.Errorf("no preserved bodies available")
	}

	value, ok := bodyContext.BodyPreserved[id]

	if !ok {
		return nil, fmt.Errorf("no preserved body found for action id %s", id)
	} else {
		return value, nil
	}
}

func (bodyContext *BodyContext) IsArray() bool {
	bodyText := string(bodyContext.CurrentBodyByteArray)
	return strings.HasPrefix(bodyText, "[") && strings.HasSuffix(bodyText, "]")
}

func (bodyContext *BodyContext) SetHeaders(headers map[string]string) {
	if headers != nil {
		if bodyContext.Headers == nil {
			bodyContext.Headers = make(map[string]string)
		}

		for key, value := range headers {
			bodyContext.Headers[key] = value
		}
	}
}

func (bodyContext *BodyContext) SetHeader(key string, value string) {
	if len(key) > 0 {
		if bodyContext.Headers == nil {
			bodyContext.Headers = make(map[string]string)
		}

		bodyContext.Headers[key] = value
	}
}

func (bodyContext *BodyContext) ParseBodyToMapObject() map[string]interface{} {
	var jsonMap map[string]interface{}
	jsonErr := json.Unmarshal(bodyContext.CurrentBodyByteArray, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) ParseBodyToMapObjectPreserved(actionId string) map[string]interface{} {
	var jsonMap map[string]interface{}
	bodyBytePreserved, _ := bodyContext.GetBodyPreserved(actionId)
	jsonErr := json.Unmarshal(bodyBytePreserved, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) ParseBodyToMapObjectArray() []map[string]interface{} {
	var jsonMap []map[string]interface{}
	jsonErr := json.Unmarshal(bodyContext.CurrentBodyByteArray, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) SetMapObject(jsonMap map[string]interface{}) {
	jsonArray, _ := json.Marshal(jsonMap)
	bodyContext.CurrentBodyByteArray = jsonArray
}

func (bodyContext *BodyContext) ConvertMapToByteArray(jsonMap map[string]interface{}) ([]byte, error) {
	jsonArray, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonArray, nil
}

func (bodyContext *BodyContext) SetArrayMapObject(arrayJsonMap []map[string]interface{}) {
	jsonArray, _ := json.Marshal(arrayJsonMap)
	bodyContext.CurrentBodyByteArray = jsonArray
}

func (bodyContext *BodyContext) GetBodyString() string {
	return string(bodyContext.CurrentBodyByteArray)
}

func (bodyContext *BodyContext) GetBody(settings *action_settings.ActionSettings) ([]byte, error) {

	if settings == nil {
		return bodyContext.CurrentBodyByteArray, nil
	}

	shouldUse, bodyRef := settings.ShouldUseBodyRef()

	if shouldUse && bodyRef != "{{bodyContext.currentBody}}" {
		bodyRef = ReplaceCalculatedValue(bodyRef)
		bodyRef = ReplacePrefixBodyContextPreserved(bodyRef)
		return bodyContext.GetBodyPreserved(bodyRef)
	} else {
		return bodyContext.CurrentBodyByteArray, nil
	}
}

func (bodyContext *BodyContext) GetBodyMap(settings *action_settings.ActionSettings) (map[string]interface{}, error) {
	bodyArray, err := bodyContext.GetBody(settings)

	if err != nil {
		return nil, err
	}

	if len(bodyArray) > 0 {
		var jsonMap map[string]interface{}
		jsonErr := json.Unmarshal(bodyArray, &jsonMap)

		if jsonErr != nil {
			return nil, jsonErr
		}
		return jsonMap, nil
	}

	return nil, nil
}

func (bodyContext *BodyContext) GetCurrentBody() []byte {
	return bodyContext.CurrentBodyByteArray
}
