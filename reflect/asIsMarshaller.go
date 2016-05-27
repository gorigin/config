package reflect

import (
	"fmt"
	"reflect"
)

// AsIsMarshaller reads data as-is, without any conversion
func AsIsMarshaller(source interface{}, target interface{}) error {
	if k := reflect.TypeOf(target).Kind(); k != reflect.Ptr {
		return fmt.Errorf("Pointer target expected, but %s provided", k)
	}

	if k := reflect.TypeOf(source).Kind(); k == reflect.Ptr {
		return fmt.Errorf("Not-pointer source expected")
	}

	sk := reflect.TypeOf(source).Kind()
	tk := reflect.TypeOf(target).Elem().Kind()
	if sk == tk {
		// Same type
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(source))
		return nil
	} else {
		return fmt.Errorf("%s != %s", sk, tk)
	}
}
