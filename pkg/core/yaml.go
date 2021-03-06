package core

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type MapStr map[string]interface{}
type Binary []byte

func (ms *MapStr) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var result map[interface{}]interface{}
	err := unmarshal(&result)
	if err != nil {
		panic(err)
	}
	*ms = cleanUpInterfaceMap(result)
	return nil
}

func (b *Binary) MarshalYAML() (interface{}, error) {
	return "base64:" + base64.StdEncoding.EncodeToString(*b), nil
}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = CleanValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = CleanValue(v)
	}
	return result
}

func cleanUpStringMap(in map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[k] = CleanValue(v)
	}
	return result
}

func CleanValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	case map[string]interface{}:
		return cleanUpStringMap(v)
	case string:
		if strings.HasPrefix(v, "base64:") {
			if decoded, err := base64.StdEncoding.DecodeString(v[7:]); err == nil {
				return Binary(decoded)
			}
		}
		return v
	case Binary:
		return v
	case int:
		return float64(v)
	case float64:
		return v
	case bool:
		return v
	case nil:
		return nil
	default:
		panic("unknown type")
	}
}
