package structutil

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func Struct2String(v any) string {
	bs, _ := json.Marshal(v)
	buf := new(bytes.Buffer)
	_ = json.Indent(buf, bs, "", "    ")
	return buf.String()
}

func Struct2Map(v any) map[string]any {
	res := make(map[string]any)

	elem := reflect.ValueOf(v).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		k := relType.Field(i).Name
		if k == "BaseModel" {
			continue
		}
		res[k] = elem.Field(i).Interface()
	}

	return res
}

func AnyIsNil(v any) bool {
	switch v.(type) {
	case string:
		return len(v.(string)) == 0
	case int:
		return v.(int) == 0
	case int8:
		return v.(int8) == 0
	case int16:
		return v.(int16) == 0
	case int32:
		return v.(int32) == 0
	case int64:
		return v.(int64) == 0
	case uint:
		return v.(uint) == 0
	case uint8:
		return v.(uint8) == 0
	case uint16:
		return v.(uint16) == 0
	case uint32:
		return v.(uint32) == 0
	case uint64:
		return v.(uint64) == 0
	case float32:
		return v.(float32) == 0
	case float64:
		return v.(float64) == 0
	case bool:
		return !v.(bool)
	default:
		return false
	}
}
