package tools

import (
	"encoding/json"
)

func SliceIntToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slice))
	for _, v := range slice {
		set[v] = struct{}{}
	}
	return set
}

func StructToMap(s interface{}) (map[string]interface{}, error) {
	structByte, err := json.Marshal(s)
	if err != nil {
		return map[string]interface{}{}, err
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(structByte, &m); err != nil {
		return map[string]interface{}{}, err
	}
	return m, err
}

func MapToStruct(m map[string]string, v interface{}) error {
	byteData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(byteData, &v); err != nil {
		return err
	}
	return nil
}
