package utils

import (
	"bytes"
	"errors"
	"testing"
)

func TestUnauthorizedMessage(t *testing.T) {
	result := UnauthorizedMessage()
	expectedResult := []byte(`{
		"success": false,
   		"message": "You unauthorized ",
	}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Expected true, false returned")
	}
}

func TestBadRequestMessage(t *testing.T) {
	result := BadRequestMessage()
	expectedResult := []byte(`{
		"success": false,
   		"message": "Bad request",
	}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Expected true, false returned")
	}
}

func TestSuccessMessage(t *testing.T) {
	result := SuccessMessage()
	expectedResult := []byte(`{
		"success": true,
   		"message": "Success",
	}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Expected true, false returned")
	}
}

func TestErrorMessage(t *testing.T) {
	err := errors.New("fail")
	result := ErrorMessage(err)
	expectedResult := []byte(`{
		"success": false,
   		"message":` + err.Error() + `,
	}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Expected true, false returned")
	}
}

func TestForbiddenMessage(t *testing.T) {
	result := ForbiddenMessage()
	expectedResult := []byte(`{
		"success": false,
   		"message": "You dont have access",
	}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Expected true, false returned")
	}
}

func TestNotFoundMessage(t *testing.T) {
	result := NotFoundMessage()
	expectedResult := []byte(`{
		"success": false,
   		"message": "Not found",
	}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Expected true, false returned")
	}
}
