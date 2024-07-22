package bencode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"torrent/pkg/assert"
)

func Unmarshal(bencoded []byte, v interface{}) error {
	decoded, remaining := Decode(bencoded)
	assert.Assert(len(remaining) == 0, "remaining data from decoding exists")
	assert.Assert(reflect.TypeOf(v).Kind() == reflect.Ptr, "attempt to unmarshal non-pointer")
	return populateStruct(decoded, reflect.Indirect(reflect.ValueOf(v)))
}

func populateStruct(data any, val reflect.Value) error {
	assert.Assert(val.Kind() == reflect.Struct, fmt.Sprintf("value must be a pointer to a struct %s %s",
		val.Kind()))
	if dict, ok := data.(map[string]interface{}); ok {
		for k, v := range dict {
			field := findFieldByTag(val, k)
			if !field.IsValid() {
				continue
			}
			assert.Assert(field.CanSet(), fmt.Sprintf("cannot set field: %s", v))
			fieldValue := reflect.ValueOf(v)
			// checks if a map is trying to unmarshall into a struct or just a normal field set
			if field.Kind() == reflect.Struct && reflect.ValueOf(v).Kind() == reflect.Map {
				assert.NoError(populateStruct(v, field), "attempt to populate map field")
			} else {
				assert.Assert(fieldValue.Type().AssignableTo(field.Type()), fmt.Sprintf("cannot assign %s to %s", fieldValue.Type(), field.Type(), field.Kind()))
				field.Set(fieldValue)
			}
		}
		return nil
	}
	return errors.New("decoded data is not a dictionary")
}

func findFieldByTag(val reflect.Value, tag string) reflect.Value {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.Tag.Get("bencode") == tag {
			return val.Field(i)
		}
	}
	return reflect.Value{}
}

func Decode(bencoded []byte) (interface{}, []byte) {
	if len(bencoded) == 0 {
		return nil, nil
	}
	switch bencoded[0] {
	case 'i':
		return decodeInteger(bencoded)
	case 'l':
		return decodeList(bencoded)
	case 'd':
		return decodeDictionary(bencoded)
	default:
		if isDigit(bencoded[0]) {
			return decodeString(bencoded)
		}
	}
	return nil, nil
}

func decodeInteger(bencoded []byte) (int, []byte) {
	end := 1
	for bencoded[end] != 'e' {
		end++
	}
	value, err := strconv.Atoi(string(bencoded[1:end]))
	if err != nil {
		return 0, nil
	}
	return value, bencoded[end+1:]
}

func decodeString(bencoded []byte) (string, []byte) {
	sep := 0
	for bencoded[sep] != ':' {
		sep++
	}
	length, err := strconv.Atoi(string(bencoded[:sep]))
	if err != nil {
		return "", nil
	}
	start := sep + 1
	end := start + length
	return string(bencoded[start:end]), bencoded[end:]
}

func decodeList(bencoded []byte) ([]interface{}, []byte) {
	var list []interface{}
	remaining := bencoded[1:]
	for remaining[0] != 'e' {
		var value interface{}
		value, remaining = Decode(remaining)
		list = append(list, value)
	}
	return list, remaining[1:]
}

func decodeDictionary(bencoded []byte) (map[string]interface{}, []byte) {
	dictionary := make(map[string]interface{})
	remaining := bencoded[1:]
	for remaining[0] != 'e' {
		var key string
		key, remaining = decodeString(remaining)
		var value interface{}
		value, remaining = Decode(remaining)
		dictionary[key] = value
	}
	return dictionary, remaining[1:]
}
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func Encode(bencoded interface{}) string {
	return ""
}
