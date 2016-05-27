package reflect

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FromStringMarshaller fills values (scalars only) using string source
func FromStringMarshaller(source string, target interface{}) error {
	if k := reflect.TypeOf(target).Kind(); k != reflect.Ptr {
		return fmt.Errorf("Pointer target expected, but %s provided", k)
	}

	kind := reflect.TypeOf(target).Elem().Kind()
	val := reflect.ValueOf(target).Elem()

	switch kind {
	case reflect.String:
		val.SetString(source)
	case reflect.Int, reflect.Int64:
		iv, err := strconv.Atoi(source)
		if err != nil {
			return err
		}
		val.SetInt(int64(iv))
	case reflect.Float64:
		fv, err := strconv.ParseFloat(source, 64)
		if err != nil {
			return err
		}
		val.SetFloat(fv)
	case reflect.Bool:
		source = strings.ToUpper(source)
		val.SetBool(source == "TRUE" || source == "T" || source == "1" || source == "YES" || source == "Y" || source == "ON")
	default:
		return fmt.Errorf("Unsupported kind %s", kind)
	}

	return nil
}
