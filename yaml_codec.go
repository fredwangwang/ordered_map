package orderedmap

import (
	"gopkg.in/yaml.v2"
)

//func iterateMap(om *OrderedMap) yaml.MapSlice {
//	ret := yaml.MapSlice{}
//
//	iter := om.IterFunc()
//	for kv, ok := iter(); ok; kv, ok = iter() {
//		kv.
//	}
//
//}

func (om OrderedMap) MarshalYAML() (interface{}, error) {
	ret := yaml.MapSlice{}

	iter := om.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		ret = append(ret, yaml.MapItem{
			Key:   kv.Key,
			Value: kv.Value,
		})
	}

	return ret, nil
}

func (om *OrderedMap) UnmarshalYAML(unmarshal func(interface{}) error) error {

	omSlice := yaml.MapSlice{}
	err := unmarshal(&omSlice)

	if err != nil {
		return err
	}

	for _, item := range omSlice {
		switch valueTyped := item.Value.(type) {
		case yaml.MapSlice:
			om.Set(item.Key, parseYAMLMap(valueTyped))
		case []interface{}:
			om.Set(item.Key, parseYAMLSlice(valueTyped))
		default:
			om.Set(item.Key, valueTyped)
		}
	}

	return nil
}

func parseYAMLSlice(content []interface{}) ([]interface{}) {
	ret := []interface{}{}

	for _, item := range content {
		switch itemTyped := item.(type) {
		case yaml.MapSlice:
			ret = append(ret, parseYAMLMap(itemTyped))
		case []interface{}:
			ret = append(ret, parseYAMLSlice(itemTyped))
		default:
			ret = append(ret, itemTyped)
		}
	}

	return ret
}

func parseYAMLMap(content yaml.MapSlice) OrderedMap {
	ret := *New()

	for _, item := range content {
		switch valueTyped := item.Value.(type) {
		case yaml.MapSlice:
			ret.Set(item.Key, parseYAMLMap(valueTyped))
		case []interface{}:
			ret.Set(item.Key, parseYAMLSlice(valueTyped))
		default:
			ret.Set(item.Key, valueTyped)
		}
	}

	return ret
}
