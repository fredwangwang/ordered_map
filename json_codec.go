package orderedmap

import (
	"encoding/json"
	"strings"
	"sort"
	"fmt"
	"container/list"
	"errors"
)

func findClosingBraces(str string, left byte, right byte) int {
	mark := 1
	isLiteral := false
	i := 1

	for ; i < len(str); i++ {
		if str[i] == '\\' {
			// consume the next symbol
			i++
		} else if str[i] == '"' {
			isLiteral = !isLiteral
		} else if !isLiteral {
			if str[i] == left {
				mark++
			} else if str[i] == right {
				mark--
			}
		}
		if mark == 0 {
			break
		}
	}
	return i
}

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
	// json require keys are strings, so hard code the key type here
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	s := string(b)
	mapToOrderedMap(om, s, m)
	return nil
}

func mapToOrderedMap(om *OrderedMap, s string, m map[string]interface{}) {
	orderedKeys := KeyIndices{}

	genericMap := map[string]interface{}{}

	// get all the keys sorted out first
	for k, _ := range m {
		kEscaped := strings.Replace(k, `"`, `\"`, -1)
		kQuoted := `"` + kEscaped + `"`
		sTrimmed := s
		for len(sTrimmed) > 0 {
			lastIndex := strings.LastIndex(sTrimmed, kQuoted)
			if lastIndex == -1 {
				break
			}
			sTrimmed = sTrimmed[0:lastIndex]
			sTrimmed = strings.TrimRight(sTrimmed, ", \n\r\t")
			maybeValidJson := sTrimmed + "}"

			// If we can successfully unmarshal the previous part, it means the match is a top-level key
			//TODO: maybe optimize here as well?
			err := json.Unmarshal([]byte(maybeValidJson), &genericMap)
			if err == nil {
				// record the position of this key in s
				ke := KeyElement{
					Key:   k,
					Index: lastIndex,
				}
				orderedKeys = append(orderedKeys, ke)
				break
			}
		}
	}
	orderedKeys = append(orderedKeys, KeyElement{Key: "", Index: len(s) - 1})
	sort.Sort(orderedKeys)

	for i := 0; i < len(orderedKeys)-1; i++ {
		contentKey := orderedKeys[i].Key
		contentKeyEscaped := `"` + strings.Replace(contentKey, `"`, `\"`, -1) + `"`
		contentEnd := orderedKeys[i+1].Index
		contentStart := orderedKeys[i].Index + len(contentKeyEscaped)
		contentStr := strings.Trim(s[contentStart:contentEnd], " \n\r:,")

		switch contentTyped := m[contentKey].(type) {
		case map[string]interface{}:
			oo := *New()
			mapToOrderedMap(&oo, contentStr, contentTyped)
			m[contentKey] = oo
		case []interface{}:
			parseSliceInMap(om, contentStr, contentTyped)
		}

	}

	li := list.List{}
	for _, ki := range orderedKeys {
		if ki.Key != "" {
			li.PushBack(ki.Key)
			om.mapper[ki.Key] = li.Back()
		}
	}

	mInter := map[interface{}]interface{}{}
	for k, v := range m {
		mInter[k] = v
	}

	om.store = mInter
}

func parseSliceInMap(om *OrderedMap, str string, content []interface{}) {
	for i, item := range content {
		switch itemTyped := item.(type) {
		case map[string]interface{}: // map
			oo := *New()
			str = str[strings.IndexByte(str, '{'):]
			idx := findClosingBraces(str, '{', '}') + 1
			mapToOrderedMap(&oo, str[:idx], itemTyped)
			content[i] = oo
			str = str[idx:]
		case []interface{}: // slice
			str = str[strings.IndexByte(str, '['):]
			idx := findClosingBraces(str, '[', ']') + 1
			parseSliceInMap(om, str[:idx], itemTyped)
			str = str[idx:]
		default: // scalar
			itemStr := fmt.Sprint(itemTyped)
			itemIdx := strings.Index(str, itemStr)
			str = str[itemIdx+len(itemStr)+1:]
		}
	}
}
