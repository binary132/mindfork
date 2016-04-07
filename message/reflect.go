package message

import (
	"fmt"
	"reflect"
)

// ReflectSet sets the target Message's value to v using reflect.
func ReflectSet(target Message, v Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				e = err
				return
			}

			panic(r)
		}
	}()

	rT := reflect.ValueOf(target)

	el := rT.Elem()
	if !el.CanSet() {
		return fmt.Errorf("unsettable target %v of type %T", e, e)
	}

	el.Set(reflect.ValueOf(v))

	return nil
}
