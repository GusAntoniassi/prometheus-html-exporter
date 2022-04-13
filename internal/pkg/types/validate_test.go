package types

import (
	"strings"
	"testing"
)

type testStruct struct {
	notRequired            string
	required               string `required`
	alsoRequired           int    `required`
	anotherStructTag       string `another-tag:"with-value"`
	RequiredWithAnotherTag string `required another-tag`
}

type complexTestStruct struct {
	notRequired string
	required    string `required`
	innerStruct testStruct
}

func TestValidateOk(t *testing.T) {
	s := testStruct{
		required:     "foo",
		alsoRequired: 0,
	}

	err := Validate(s)
	if err != nil {
		t.Fatalf("validating with required fields set should not return any errors. error returned was: %s", err.Error())
	}
}

func TestValidateWithNestedStruct(t *testing.T) {
	s := complexTestStruct{
		required: "foo",
		innerStruct: testStruct{
			required:     "foo",
			alsoRequired: 0,
		},
	}

	err := Validate(s)
	if err != nil {
		t.Fatalf("validating with required fields set should not return any errors. error returned was: %s", err.Error())
	}
}

func TestValidateWithErrors(t *testing.T) {
	s := testStruct{
		notRequired:            "foobar",
		anotherStructTag:       "foobar",
		RequiredWithAnotherTag: "foobar",
	}

	err := Validate(s)
	if err == nil {
		t.Fatal("validating without required fields set should return an error")
	}

	errorMessage := err.Error()
	if !strings.Contains(errorMessage, "field required") || !strings.Contains(errorMessage, "field alsoRequired") {
		t.Fatal("expected error message to contain all invalid fields, got the following message:\n", errorMessage)
	}
}
