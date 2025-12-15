package opera

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Portions of this file are derived from github.com/samber/mo, licensed under the MIT License.

// MarshalJSON encodes Option into json.
// Go 1.20+ relies on the IsZero method when the `omitempty` tag is used
// unless a custom MarshalJSON method is defined.  Then the IsZero method is ignored.
// current best workaround is to instead use `omitzero` tag with Go 1.24+
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.hasVal {
		return json.Marshal(o.val)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON decodes Option from json.
func (o *Option[T]) UnmarshalJSON(b []byte) error {
	o.val = Empty[T]() // reset the value if not set later.

	// If user manually set the field to be `null`, then it either means the option is absent or present with a zero value.
	if bytes.Equal([]byte("null"), bytes.ToLower(b)) {
		// // If the type is a pointer, then it means the option is present with a zero value.
		// o.isPresent = reflect.TypeOf(o.value).Kind() == reflect.Ptr
		// return nil

		o.hasVal = false
		return nil
	}

	err := json.Unmarshal(b, &o.val)
	if err != nil {
		return err
	}

	o.hasVal = true
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (o Option[T]) MarshalText() ([]byte, error) {
	return json.Marshal(o)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (o *Option[T]) UnmarshalText(data []byte) error {
	return json.Unmarshal(data, o)
}

// MarshalBinary is the interface implemented by an object that can marshal itself into a binary form.
func (o Option[T]) MarshalBinary() ([]byte, error) {
	if !o.hasVal {
		return []byte{0}, nil
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(o.val); err != nil {
		return []byte{}, err
	}

	return append([]byte{1}, buf.Bytes()...), nil
}

// UnmarshalBinary is the interface implemented by an object that can unmarshal a binary representation of itself.
func (o *Option[T]) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("Option[T].UnmarshalBinary: no data")
	}

	if data[0] == 0 {
		o.hasVal = false
		o.val = Empty[T]()
		return nil
	}

	buf := bytes.NewBuffer(data[1:])
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&o.val)
	if err != nil {
		return err
	}

	o.hasVal = true
	return nil
}

// GobEncode implements the gob.GobEncoder interface.
func (o Option[T]) GobEncode() ([]byte, error) {
	return o.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (o *Option[T]) GobDecode(data []byte) error {
	return o.UnmarshalBinary(data)
}

// Scan implements the SQL sql.Scanner interface.
func (o *Option[T]) Scan(src any) error {
	if src == nil {
		o.hasVal = false
		o.val = Empty[T]()
		return nil
	}

	// is is only possible to assert interfaces, so convert first
	// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#why-not-permit-type-assertions-on-values-whose-type-is-a-type-parameter
	var t T
	if tScanner, ok := interface{}(&t).(sql.Scanner); ok {
		if err := tScanner.Scan(src); err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}

		o.hasVal = true
		o.val = t
		return nil
	}

	if av, err := driver.DefaultParameterConverter.ConvertValue(src); err == nil {
		if v, ok := av.(T); ok {
			o.hasVal = true
			o.val = v
			return nil
		}
	}

	return o.scanConvertValue(src)
}

// Value implements the driver Valuer interface.
func (o Option[T]) Value() (driver.Value, error) {
	if !o.hasVal {
		return nil, nil
	}

	return driver.DefaultParameterConverter.ConvertValue(o.val)
}

// Equal compares two Option[T] instances for equality
func (o Option[T]) Equal(other Option[T]) bool {
	if !o.hasVal && !other.hasVal {
		return true
	}

	if o.hasVal != other.hasVal {
		return false
	}

	return reflect.DeepEqual(o.val, other.val)
}

// scanConvertValue tries to scan src into Option[T] when T does not implement sql.Scanner.
func (o *Option[T]) scanConvertValue(src any) error {
	var st sql.Null[T]
	if err := st.Scan(src); err == nil {
		o.hasVal = true
		o.val = st.V
		return nil
	}
	return fmt.Errorf("failed to scan Option[T]")
}
