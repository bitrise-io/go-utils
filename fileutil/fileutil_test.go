package fileutil

import (
	"testing"
)

func TestJSONUnmarshalFromFile(t *testing.T) {
	testJSONString := `{
	"this_will": "it was decoded",
	"this_wont": "but this was not"
}`

	type MyType struct {
		ThisWill string `json:"this_will"`
		thisWont string `json:"this_wont"`
	}

	var test MyType
	if err := JSONUnmarshalFromBytes([]byte(testJSONString), &test); err != nil {
		t.Fatal(err)
	}
	if test.ThisWill != "it was decoded" {
		t.Error("ThisWill was not correctly decoded (it should!)")
	}
	if test.thisWont != "" {
		t.Error("thisWont was decoded??")
	}
}

func TestJSONMarshall(t *testing.T) {
	testJSONString := `{"this_will":"it was decoded"}`

	type MyType struct {
		ThisWill string `json:"this_will"`
		thisWont string `json:"this_wont"`
	}
	test := MyType{
		ThisWill: "it was decoded",
		thisWont: "but this was not",
	}
	bytes, err := JSONMarshall(test, false)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != testJSONString {
		t.Errorf("Not correctly marshaled\n (%s)\n should be\n (%s)", string(bytes), testJSONString)
	}

	testJSONString = `{
	"this_will": "it was decoded"
}`
	bytes, err = JSONMarshall(test, true)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != testJSONString {
		t.Errorf("\nNot correctly marshaled\n%s\nShould be\n%s", string(bytes), testJSONString)
	}
}
