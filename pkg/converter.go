package app

import (
	"errors"
	"reflect"
	"strconv"
)

// MakeArrFromStruct struct -> []interface
func MakeArrFromStruct(data interface{}) []interface{} {
	arr := []interface{}{}
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		arr = append(arr, v.Field(i).Interface())
	}
	return arr
}

// MakeMapFromStruct struct -> map[string]interface
func MakeMapFromStruct(data interface{}) map[string]interface{} {
	if reflect.TypeOf(data).Kind() == reflect.Map {
		return data.(map[string]interface{})
	}

	structLen := reflect.ValueOf(data).NumField()
	t := reflect.TypeOf(data)
	result := map[string]interface{}{}

	for i := 0; i < structLen; i++ {
		result[t.Field(i).Name] = ""
		reflect.ValueOf(result[t.Field(i).Name]).Set(reflect.ValueOf(t.Field(i)))
	}

	return result
}

// FillStructFromArr fill the current struct from array by order
func FillStructFromArr(sampleStruct interface{}, data []interface{}) error {
	structValue := reflect.ValueOf(sampleStruct).Elem()

	for i, v := range data {
		structFieldValue := structValue.FieldByIndex([]int{i})

		if !structFieldValue.IsValid() {
			return errors.New("no such field index: " + strconv.Itoa(i) + " in obj")
		}
		if !structFieldValue.CanSet() {
			return errors.New("can't set " + strconv.Itoa(i) + " field index")
		}

		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(v)

		if structFieldType != val.Type() {
			return errors.New("provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}
	return nil
}
