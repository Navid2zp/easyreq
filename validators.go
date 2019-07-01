package easyreq

import (
	"errors"
	"reflect"
)

// Validates the request data
func (r *Request) validate() error {
	// If no URL has been set
	if r.URL == "" {
		return errors.New("no url specified")
	}
	return nil
}

func validateSaveToInterface(saveTo interface{}) error {
	// Check if `saveTo` is a pointer
	if reflect.TypeOf(saveTo).Kind() != reflect.Ptr {
		return errors.New("non-pointer " + reflect.TypeOf(saveTo).Kind().String() + ")")
	}
	return nil
}
