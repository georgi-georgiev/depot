package tst

import (
	"bytes"
	"encoding/json"
	"strings"
)

func Prettify(input []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, input, "", "  ")
	return out.Bytes()
}

//TODO: isn't it better to get data from DB?
func Bind(expected []byte, actual []byte, criteria string) []byte {
	var actualData interface{}
	err := json.Unmarshal(actual, &actualData)
	if err != nil {
		panic(err)
	}
	actualMsg := actualData.(map[string]interface{})

	var expectedData interface{}
	err = json.Unmarshal(expected, &expectedData)
	if err != nil {
		panic(err)
	}
	expectedMsg := expectedData.(map[string]interface{})

	fields := strings.Split(criteria, ".")
	field := fields[0]

	if actualMsg[field] == nil {
		panic("missing field:" + field)
	}

	expectedMsg[field] = actualMsg[field]

	actual, err = json.Marshal(actualMsg)
	if err != nil {
		panic(err)
	}

	expected, err = json.Marshal(expectedMsg)
	if err != nil {
		panic(err)
	}
	if len(fields) > 2 {
		return Bind(expected, actual, strings.Join(fields[1:], "."))
	} else {
		return expected
	}
}

func BindArray(expected []byte, actual []byte, criteria string) []byte {
	var actualData []interface{}
	err := json.Unmarshal(actual, &actualData)
	if err != nil {
		panic(err)
	}

	var expectedData []interface{}
	err = json.Unmarshal(expected, &expectedData)
	if err != nil {
		panic(err)
	}

	var expectedArray []byte

	expectedArray = append(expectedArray, byte('['))

	for i := 0; i < len(expectedData); i++ {
		expectedMarshal, err := json.Marshal(expectedData[i])
		if err != nil {
			panic(err)
		}

		actualMarshal, err := json.Marshal(actualData[i])
		if err != nil {
			panic(err)
		}

		expectedArray = append(expectedArray, Bind(expectedMarshal, actualMarshal, criteria)...)
	}

	expectedArray = append(expectedArray, byte(']'))

	return expectedArray
}
