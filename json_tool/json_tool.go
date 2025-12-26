package json_tool

import "encoding/json"

func ToJson(v any) string {

	if v == nil {
		return ""
	}
	marshal, _ := json.Marshal(v)
	return string(marshal)
}
