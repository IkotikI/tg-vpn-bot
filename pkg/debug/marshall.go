package debug

import "encoding/json"

func JSON(v any) string {

	json, err := json.MarshalIndent(v, "", "    ")

	if err != nil {
		return "<debug.json: invalid json>"
	}

	return string(json)

}
