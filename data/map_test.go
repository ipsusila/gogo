package data

import (
	"encoding/json"
	"fmt"
	"testing"
)

const jsonDoc = `
{
	"localhost": {
		"tag": "dev_latest",
		"vhost": "localhost.com"
	},
	"development": {
		"tag": "dev_latest",
		"vhost": "dev.com"
	},
	"other": 123,
	"release": {
		"DB": {
			"host": "localhost",
			"port": "5432"
		}
	}
}
`

func TestFlatMap(t *testing.T) {
	//JSON to map
	var m1 map[string]interface{}
	err := json.Unmarshal([]byte(jsonDoc), &m1)
	if err != nil {
		t.Fatal(err)
	}
	//Flag test
	if IsFlatMap(m1) {
		t.FailNow()
	}
	fm1 := FlatMap(m1, ".")
	fmt.Printf("%#v\n", fm1)

	//Custom map
	m2s := make(map[interface{}]interface{})
	m2s[456] = "Num: 456"
	m2s[true] = "Num: true"

	m2 := make(map[interface{}]interface{})
	m2["id"] = "Flat_map_test"
	m2["sub"] = m2s
	m2[8080] = "HTTP-PORT"
	if IsFlatMap(m2) {
		t.FailNow()
	}
	fm2 := FlatMap(m2, ".")
	fmt.Printf("%#v\n", fm2)
}
