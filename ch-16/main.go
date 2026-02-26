package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"chap_16/internal/orders"
	"unsafe"
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

func printSizesAndOffsets(inp any) {
	switch v := inp.(type) {
	case orders.OrderInfo:
		fmt.Println("OI, Total Size:", unsafe.Sizeof(v))
		fmt.Printf("Offsets: %v field 1, %v field 2, %v field 3, %v field 4, %v field 5\n",
			unsafe.Offsetof(v.OrderCode),
			unsafe.Offsetof(v.Amount),
			unsafe.Offsetof(v.OrderNumber),
			unsafe.Offsetof(v.Items),
			unsafe.Offsetof(v.IsReady),
		)
	case orders.SmallOrderInfo:
		fmt.Println("SOI, Total Size:", unsafe.Sizeof(v))
		fmt.Printf("Offsets: %v field 1, %v field 2, %v field 3, %v field 4, %v field 5\n",
			unsafe.Offsetof(v.Items),
			unsafe.Offsetof(v.Amount),
			unsafe.Offsetof(v.OrderCode),
			unsafe.Offsetof(v.OrderNumber),
			unsafe.Offsetof(v.IsReady),
		)
	}
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
	// 2
	defaultOI := orders.OrderInfo{}
	smallOI := orders.SmallOrderInfo{}
	printSizesAndOffsets(defaultOI)
	printSizesAndOffsets(smallOI)
}
