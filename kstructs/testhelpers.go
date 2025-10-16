// Package kstructs is deprecated: use ksqltest instead.
package kstructs

import (
	"fmt"
	"reflect"

	"github.com/vingarcia/kbuilder/helpers/structs"
)

// StructToMap converts any struct type to a map based on
// the tag named `ksql`, i.e. `ksql:"map_key_name"`
//
// Valid pointers are dereferenced and copied to the map,
// null pointers are ignored.
//
// This function is efficient in the fact that it caches
// the slower steps of the reflection required to perform
// this task.
//
// deprecated: Use ksqltest.StructToMap() instead
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	return structs.StructToMap(obj)
}

// FillStructWith is meant to be used on unit tests to mock
// the response from the database.
//
// The first argument is any struct you are passing to a ksql func,
// and the second is a map representing a database row you want
// to use to update this struct.
//
// deprecated: Use ksqltest.FillStructWith() instead
func FillStructWith(record interface{}, dbRow map[string]interface{}) error {
	v := reflect.ValueOf(record)
	t := v.Type()

	if t.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"FillStructWith: expected input to be a pointer to struct but got %T",
			record,
		)
	}

	t = t.Elem()
	v = v.Elem()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf(
			"FillStructWith: expected input to be a pointer to a struct, but got %T",
			record,
		)
	}

	info, err := structs.GetTagInfo(t)
	if err != nil {
		return err
	}

	for colName, rawSrc := range dbRow {
		fieldInfo := info.ByName(colName)
		if !fieldInfo.Valid {
			// Ignore columns not tagged with `ksql:"..."`
			continue
		}

		src := structs.NewPtrConverter(rawSrc)
		dest := v.Field(fieldInfo.Index)
		destType := t.Field(fieldInfo.Index).Type

		destValue, err := src.Convert(destType)
		if err != nil {
			return fmt.Errorf("FillStructWith: error on field `%s`: %w", colName, err)
		}

		dest.Set(destValue)
	}

	return nil
}

// FillSliceWith is meant to be used on unit tests to mock
// the response from the database.
//
// The first argument is any slice of structs you are passing to a ksql func,
// and the second is a slice of maps representing the database rows you want
// to use to update this struct.
//
// deprecated: Use ksqltest.FillSliceWith() instead
func FillSliceWith(entities interface{}, dbRows []map[string]interface{}) error {
	sliceRef := reflect.ValueOf(entities)
	sliceType := sliceRef.Type()
	if sliceType.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"FillSliceWith: expected input to be a pointer to a slice of structs but got %v",
			sliceType,
		)
	}

	structType, isSliceOfPtrs, err := structs.DecodeAsSliceOfStructs(sliceType.Elem())
	if err != nil {
		return fmt.Errorf("FillSliceWith: %w", err)
	}

	slice := sliceRef.Elem()
	for idx, row := range dbRows {
		if slice.Len() <= idx {
			var elemValue reflect.Value
			elemValue = reflect.New(structType)
			if !isSliceOfPtrs {
				elemValue = elemValue.Elem()
			}
			slice = reflect.Append(slice, elemValue)
		}

		err := FillStructWith(slice.Index(idx).Addr().Interface(), row)
		if err != nil {
			return fmt.Errorf("FillSliceWith: %w", err)
		}
	}

	sliceRef.Elem().Set(slice)

	return nil
}

// CallFunctionWithRows was created for helping test the QueryChunks method
//
// deprecated: Use ksqltest.CallFunctionWithRows() instead
func CallFunctionWithRows(fn interface{}, rows []map[string]interface{}) error {
	fnValue := reflect.ValueOf(fn)
	chunkType, err := structs.ParseInputFunc(fn)
	if err != nil {
		return err
	}

	chunk := reflect.MakeSlice(chunkType, 0, len(rows))

	// Create a pointer to a slice (required by FillSliceWith)
	chunkPtr := reflect.New(chunkType)
	chunkPtr.Elem().Set(chunk)

	err = FillSliceWith(chunkPtr.Interface(), rows)
	if err != nil {
		return err
	}

	err, _ = fnValue.Call([]reflect.Value{chunkPtr.Elem()})[0].Interface().(error)

	return err
}
