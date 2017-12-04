package gojsonhook

import (
	"encoding/json"
	"reflect"
)

// PreMarshaler is the interface that can be implemented if you want to do something before marshaling to json
type PreMarshaler interface {
	BeforeMarshal() error
}

// PostUnmarshaler is the interface that can be implemented if you want to do something after unmarshaling to some type
// AfterUnmarshal is called after json.Unmarshal, so if there is an error in AfterMarshal, the input value may be modified,
// so the input value is not credible if any error has happened
type PostUnmarshaler interface {
	AfterUnmarshal() error
}

var (
	preMarshalType    = reflect.TypeOf(new(PreMarshaler)).Elem()
	postUnmarshalType = reflect.TypeOf(new(PostUnmarshaler)).Elem()
)

// Marshal Call BeforeMarshal if defined and marshal the input interface
func Marshal(v interface{}) ([]byte, error) {
	vType := reflect.TypeOf(v)

	if vType.Implements(preMarshalType) {
		err := v.(PreMarshaler).BeforeMarshal()
		if err != nil {
			return nil, err
		}
	}

	if vType.Kind() != reflect.Ptr {
		if reflect.PtrTo(vType).Implements(preMarshalType) {
			ptr := reflect.New(vType)
			temp := ptr.Elem()
			temp.Set(reflect.ValueOf(v))

			ptr.MethodByName("BeforeMarshal").Call(nil)
			v = temp.Interface()
		}
	}

	bytes, err := json.Marshal(v)
	return bytes, err

}

// Unmarshal Unmarshal data and try to call AfterUnmarshal if defined
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	vType := reflect.TypeOf(v)
	if vType.Implements(postUnmarshalType) {
		err := v.(PostUnmarshaler).AfterUnmarshal()
		if err != nil {
			return err
		}
	}

	return nil
}
