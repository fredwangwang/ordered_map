package orderedmap

import (
	"testing"
	"fmt"
	"gopkg.in/yaml.v2"
)

func TestMarshalYAML(t *testing.T) {
	o := New()
	// number
	o.Set("number", 3)
	// string
	o.Set("string", "x")
	// new value keeps key in old position
	o.Set("number", 4)
	// keys not sorted alphabetically
	o.Set("z", 1)
	o.Set("a", 2)
	o.Set("b", 3)
	// slice
	o.Set("slice", []interface{}{
		"1",
		1,
	})
	// orderedmap
	v := New()
	v.Set("e", 1)
	v.Set("a", 2)
	o.Set("orderedmap", v)
	// double quote in key
	o.Set(`test"ing`, 9)
	bi, err := yaml.Marshal(o)
	if err != nil {
		t.Error("Marshalling yaml", err)
	}
	si := string(bi)
	ei := `number: 4
string: x
z: 1
a: 2
b: 3
slice:
- "1"
- 1
orderedmap:
  e: 1
  a: 2
test"ing: 9
`
	if si != ei {
		fmt.Println(ei)
		fmt.Println(si)
		t.Error("YAML MarshalIndent value is incorrect", si)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	s := `---
number: 4
string: x
z: 1
a: should not break with unclosed { character in value
b: 3
slice:
- '1'
- 1
orderedmap:
  e: 1
  a { nested key with brace: with a }}}} }} {{{ brace value
  after:
    link: test {{{ with even deeper nested braces }
test"ing: 9
after: 1
multitype_array:
- test
- 1
- map: obj
  it: 5
  ":colon in key": 'colon: in value'
- - inner: map
should not break with { character in key: 1
`
	om := New()
	err := yaml.Unmarshal([]byte(s), &om)
	if err != nil {
		t.Error("YAML Unmarshal error", err)
	}
	// Check the root keys
	expectedKeys := []string{
		"number",
		"string",
		"z",
		"a",
		"b",
		"slice",
		"orderedmap",
		"test\"ing",
		"after",
		"multitype_array",
		"should not break with { character in key",
	}

	iter := om.IterFunc()
	i := 0
	for kv, ok := iter(); ok; kv, ok = iter() {
		if kv.Key != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
		}
		i++
	}
	// Check nested maps are converted to orderedmaps
	// nested 1 level deep
	expectedKeys = []string{
		"e",
		"a { nested key with brace",
		"after",
	}
	vi, ok := om.Get("orderedmap")
	if !ok {
		t.Error("Missing key for nested map 1 deep")
	}
	v := vi.(OrderedMap)
	iter = v.IterFunc()
	i = 0
	for kv, ok := iter(); ok; kv, ok = iter() {
		if kv.Key != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
		}
		i++
	}
	// nested 2 levels deep
	expectedKeys = []string{
		"link",
	}
	vi, ok = v.Get("after")
	if !ok {
		t.Error("Missing key for nested map 2 deep")
	}
	v = vi.(OrderedMap)
	iter = v.IterFunc()
	i = 0
	for kv, ok := iter(); ok; kv, ok = iter() {
		if kv.Key != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
		}
		i++
	}
	// multitype array
	expectedKeys = []string{
		"map",
		"it",
		":colon in key",
	}
	vislice, ok := om.Get("multitype_array")
	if !ok {
		t.Error("Missing key for multitype array")
	}
	vslice := vislice.([]interface{})
	vmap := vslice[2].(OrderedMap)
	iter = vmap.IterFunc()
	i = 0
	for kv, ok := iter(); ok; kv, ok = iter() {
		if kv.Key != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
		}
		i++
	}

	expectedKeys = []string{"inner"}
	vinnerslice := vslice[3].([]interface{})
	vinnermap := vinnerslice[0].(OrderedMap)
	iter = vinnermap.IterFunc()
	i = 0
	for kv, ok := iter(); ok; kv, ok = iter() {
		if kv.Key != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
		}
		i++
	}
}

func TestUnmarshalYAMLSpecialChars(t *testing.T) {
	s := `{ " \\\\\\\\\\\\ "  : { "\\\\\\" : "\\\\\"\\" }, "\\":  " \\\\ test " }`
	o := New()
	err := yaml.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Error("yaml Unmarshal error with special chars", err)
	}
}

func TestUnmarshalYAMLArrayOfMaps(t *testing.T) {
	s := `
---
name: test
percent: 6
breakdown:
- name: a
  percent: 0.9
- name: b
  percent: 0.9
- name: d
  percent: 0.4
- name: e
  percent: 2.7

`
	om := New()
	err := yaml.Unmarshal([]byte(s), &om)
	if err != nil {
		t.Error("yaml Unmarshal error", err)
	}
	// Check the root keys
	expectedKeys := []string{
		"name",
		"percent",
		"breakdown",
	}
	iter := om.IterFunc()
	i := 0
	for kv, ok := iter(); ok; kv, ok = iter() {
		if kv.Key != expectedKeys[i] {
			t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
		}
		i++
	}
	// Check nested maps are converted to orderedmaps
	// nested 1 level deep
	expectedKeys = []string{
		"name",
		"percent",
	}
	vi, ok := om.Get("breakdown")
	if !ok {
		t.Error("Missing key for nested map 1 deep")
	}
	vs := vi.([]interface{})
	for _, vInterface := range vs {
		v := vInterface.(OrderedMap)
		iter := v.IterFunc()
		i := 0
		for kv, ok := iter(); ok; kv, ok = iter() {
			if kv.Key != expectedKeys[i] {
				t.Error("Unmarshal root key order", i, kv.Key, "!=", expectedKeys[i])
			}
			i++
		}
	}
}
