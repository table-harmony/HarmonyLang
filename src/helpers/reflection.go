package helpers

import (
	"fmt"
	"reflect"
)

// Expected type is passed as a generic and this method will use reflection to compare the underlying type agains T.
// Returns the casted type or error if it fails.
func ExpectType[T any](r any) (T, error) {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	recievedType := reflect.TypeOf(r)

	if expectedType == recievedType {
		return r.(T), nil
	}

	var zeroValue T
	return zeroValue, fmt.Errorf("expected %v but instead received %v", expectedType, recievedType)
}
