package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type strLenStruct struct {
	FieldOne, FieldTwo string `minStrLen:"10"`
}

func ValidateStringLength(inp any) error {
	it := reflect.TypeOf(inp)
	iv := reflect.ValueOf(inp)

	if it.Kind() != reflect.Struct {
		return errors.New("Expected struct as an input")
	}

	var errWrapper error
	for i := 0; i < it.NumField(); i++ {
		curField := it.Field(i)
		if tag := curField.Tag.Get("minStrLen"); tag != "" && curField.Type.Name() == "string" {
			conv, err := strconv.Atoi(tag)
			if err != nil {
				errWrapper = errors.Join(errWrapper, err)
			}
			if strLen := len(iv.Field(i).Interface().(string)); strLen < conv {
				errWrapper = errors.Join(errWrapper, fmt.Errorf("Length doesn't match the tag - %d (string - %s): expected equals or more than %d", strLen, iv.Field(i).Interface().(string), conv))
			}
		}
	}
	return errWrapper
}

func main() { // 1
	structOne := strLenStruct{"some symbols", "more symbols"}
	structTwo := strLenStruct{"invalid", ""}
	randData := "Some random datatype"
	inputs := []any{structOne, structTwo, randData}
	for i := range inputs {
		if err := ValidateStringLength(inputs[i]); err != nil {
			fmt.Printf("%v (%v) produced an error: %v\n", reflect.TypeOf(inputs[i]).Name(), inputs[i], err)
		}
	}
}
