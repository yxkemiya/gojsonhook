package gojsonhook

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

type TestCustomMarshal struct {
	CreatedAt   time.Time `json:"-"`
	CreatedTime int64     `json:"created_at"`
}

func (tcm *TestCustomMarshal) BeforeMarshal() error {
	tcm.CreatedTime = tcm.CreatedAt.Unix()
	return nil
}

func (tcm *TestCustomMarshal) AfterUnmarshal() error {
	tcm.CreatedAt = time.Unix(tcm.CreatedTime, 0)
	return nil
}

type TestCustomMarshalWithError struct {
	CreatedAt   time.Time `json:"-"`
	CreatedTime int64     `json:"created_at"`
}

var (
	ErrBeforeMarshalTest  = errors.New("test before marshal error")
	ErrAfterUnmarshalTest = errors.New("test after unmarshal error")
)

func (tcmwe *TestCustomMarshalWithError) BeforeMarshal() error {
	tcmwe.CreatedTime = tcmwe.CreatedAt.Unix()
	return ErrBeforeMarshalTest
}

func (tcmwe *TestCustomMarshalWithError) AfterUnmarshal() error {
	tcmwe.CreatedAt = time.Unix(tcmwe.CreatedTime, 0)
	return ErrAfterUnmarshalTest
}

func Test_CustomeMarshaler(t *testing.T) {
	custom := TestCustomMarshal{
		CreatedAt: time.Unix(time.Now().Unix(), 0),
	}

	customBytes1, err := Marshal(custom)
	if err != nil {
		t.Fatalf("custom marshal json error [%+v]", err)
	}
	t.Logf("%s", customBytes1)

	customBytes2, err := Marshal(&custom)
	if err != nil {
		t.Fatalf("custom marshal pointer type error [%+v]", err)
	}
	t.Logf("%s", customBytes2)

	if bytes.Compare(customBytes1, customBytes2) != 0 {
		t.Fatalf("custom marshal value not equal")
	}

	var cus TestCustomMarshal
	err = Unmarshal(customBytes2, &cus)
	if err != nil {
		t.Fatalf("custom unmarshal json error [%+v]", err)
	}
	if cus.CreatedAt != custom.CreatedAt {
		t.Fatal("custom unmarshal json has unexpected result")
	}
	t.Log(cus)

	customWithErr := &TestCustomMarshalWithError{
		CreatedAt: time.Unix(time.Now().Unix(), 0),
	}

	customErrBytes1, err := Marshal(customWithErr)
	if err != ErrBeforeMarshalTest {
		t.Fatalf("test custome marshal with error has unexpected error [%+v]", err)
	}
	if customErrBytes1 != nil {
		t.Fatalf("custom marshal with error should have nil bytes")
	}

	var cusErr TestCustomMarshalWithError
	err = Unmarshal(customBytes2, &cusErr)
	if err != ErrAfterUnmarshalTest {
		t.Fatalf("test custom unmarshal with error has unexpected error [%+v]", err)
	}
}
