# Ordered Map for golang

[![Build Status](https://travis-ci.org/fredwangwang/orderedmap.svg?branch=master)](https://travis-ci.org/fredwangwang/orderedmap)

**OrderedMap** is a Python port of OrderedDict implemented in golang. 
Golang's builtin `map` purposefully randomizes the iteration of stored key/values. 
**OrderedMap** struct preserves inserted key/value pairs; such that on iteration, 
key/value pairs are received in inserted (first in, first out) order.

## Features
- Full support Key/Value for all data types
- Exposes an Iterator that iterates in order of insertion
- Full Get/Set/Delete map interface
- Supports Golang v1.3 through v1.10
- Supports **JSON** Marshal/Unmarshal
- Supports **YAML** Marshal/Unmarshal

## Download and Install 
  
`go get github.com/fredwangwang/orderedmap`

## Examples

### Create, Get, Set, Delete

```go
package main

import (
    "fmt"
    "github.com/fredwangwang/orderedmap"
)

func main() {

    // Init new OrderedMap
    om := orderedmap.New()

    om.Set("a", 1)
    om.Set("b", 2)

    if val, ok := om.Get("b"); ok == true {
        fmt.Println(val)
    }

    om.Delete("a")

    if _, ok := om.Get("a"); ok == false {
        fmt.Println("c not found")
    }
    
    fmt.Println(om)
}
```


### Iterator

```go
n := 100
om := orderedmap.New()

for i := 0; i < n; i++ {
    om.Set(i, fmt.Sprintf("%d", i * i))
}

// Iterate though values
// - Values iteration are in insert order
// - Returned in a key/value pair struct
iter := om.IterFunc()
for kv, ok := iter(); ok; kv, ok = iter() {
    fmt.Println(kv, kv.Key, kv.Value)
}
```

### Custom Structs

```go
om := orderedmap.New()
om.Set("one", &MyStruct{1, 1.1})
om.Set("two", &MyStruct{2, 2.2})
om.Set("three", &MyStruct{3, 3.3})

fmt.Println(om)
// Ouput: OrderedMap[one:&{1 1.1},  two:&{2 2.2},  three:&{3 3.3}, ]
```

### JSON marshal & unmarshal
```go
rawPayload := `{
"number": 4,
"string": "x",
"z": 1,
"a": 2,
"b": 3,
"slice": [
  "1",
  1
],
"orderedmap": {
  "e": 1,
  "a": 2
},
"test\"ing": 9
}`

om := orderedmap.New()
json.Unmarshal([]byte(rawPayload), om)
constructedPayload, _ := json.MarshalIndent(om, "", "  ")
println(string(constructedPayload)) // will get you the same thing as rawPayload
```

### YAML marshal & unmarshal
```go
rawPayload := `number: 4
string: x
z: 1
a: 2
b: 3
slice:
- '1'
- 1
orderedmap:
  e: 1
  a: 2
testing: 9`

om := orderedmap.New()
yaml.Unmarshal([]byte(rawPayload), om)
constructedPayload, _ := yaml.Marshal(om)
println(string(constructedPayload)) // will get you the same thing as rawPayload
```

  
## For Development

Git clone project 

`git clone https://github.com/fredwangwang/orderedmap.git`  
  
Build and install project

`make`

Run tests 

`make test`







