package lib

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

func GenerateRandom(k int) string {
	const CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	n := int64(len(CHARS))

	code := ""
	for i := 1; i <= k; i++ {
		position, _ := rand.Int(rand.Reader, big.NewInt(n))
		code += string(CHARS[position.Int64()])
	}

	return code
}

func GenerateKey(size int) []byte {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	return key
}

func GenerateID(prefix string) string {
	now := time.Now().Unix()
	random := GenerateRandom(5)

	return fmt.Sprintf("%s_%x_%s", prefix, now, random)
}

func CompareStructs(obj1, obj2 interface{}) bool {
	value1 := reflect.ValueOf(obj1)
	value2 := reflect.ValueOf(obj2)

	// Check if both values are structs
	if value1.Kind() != reflect.Struct || value2.Kind() != reflect.Struct {
		return false
	}

	// Get the type of the structs
	type1 := value1.Type()
	type2 := value2.Type()

	// Check if the structs have the same number of fields
	if type1.NumField() != type2.NumField() {
		return false
	}

	// Iterate over the fields and compare their values
	for i := 0; i < type1.NumField(); i++ {
		field1 := value1.Field(i)
		field2 := value2.Field(i)

		if field1.Type() == reflect.TypeOf(time.Time{}) {
			// Handle time.Time fields separately
			time1 := field1.Interface().(time.Time)
			time2 := field2.Interface().(time.Time)

			// Compare the Unix timestamps of the time values
			if time1.Unix() != time2.Unix() {
				return false
			}
		} else if field1.Kind() == reflect.Struct ||
			field1.Kind() == reflect.Map ||
			field1.Kind() == reflect.Slice ||
			field2.Kind() == reflect.Struct ||
			field2.Kind() == reflect.Map ||
			field2.Kind() == reflect.Slice {
			// Recursively check structs
			ok := CompareStructs(field1.Interface(), field2.Interface())
			if !ok {
				return false
			}
		} else {
			// Compare other field values
			if !reflect.DeepEqual(field1.Interface(), field2.Interface()) {
				return false
			}
		}
	}

	return true
}

func MapToStruct(m map[string]interface{}, target interface{}) error {
	// Check if target is a non-nil struct
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return errors.New("lib.MapToStruct(): target must be a non-nil pointer")
	}

	targetElem := targetValue.Elem()
	if targetElem.Kind() != reflect.Struct {
		return errors.New("lib.MapToStruct(): target must be a pointer to a struct")
	}

	// Convert map data to json bytes
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("lib.MapToStruct(): %w", err)
	}

	// Convert json bytes to target struct
	err = json.Unmarshal(b, &target)
	if err != nil {
		return fmt.Errorf("lib.MapToStruct(): %w", err)
	}

	return nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	// Check if obj is non-nil and struct
	if obj == nil {
		return nil, errors.New("lib.StructToMap(): Object must be non-nil and struct")
	}

	v := reflect.TypeOf(obj)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("lib.StructToMap(): target must be a struct")
	}

	// Convert obj to json bytes
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("lib.StructToMap(): %w", err)
	}

	// Convert json bytes to map[string]interface{}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(b, &mapData)
	if err != nil {
		return nil, fmt.Errorf("lib.StructToMap(): %w", err)
	}

	return mapData, nil
}

// Remove leading and trailing whitespaces
func Trim(text string) string {
	whitespaces := []rune{' ', '　', '\n', '\t'}
	text = strings.TrimFunc(text, func(r rune) bool {
		for _, c := range whitespaces {
			if c == r {
				return true
			}
		}
		return false
	})

	return text
}
