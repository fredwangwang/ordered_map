package orderedmap

import (
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

func (om OrderedMap) MarshalJSON() ([]byte, error) {
	// as JSON specifies key must be string, check if the keys can be casted to string first
	iter := om.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		_, ok := kv.Key.(string)
		if ok != true {
			return nil, errors.New(fmt.Sprintf("key: %v is not string", kv.Key))
		}
	}

	s := "{"
	isNext := false
	iter = om.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		if isNext {
			s += ","
		}
		k := kv.Key.(string)
		s += fmt.Sprintf("\"%s\":", strings.Replace(k, `"`, `\"`, -1))

		v := kv.Value
		vBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		s += string(vBytes)
		isNext = true
	}
	s += "}"
	return []byte(s), nil
}

func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	return yaml.Unmarshal(b, om)
}
