package model

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type Field struct {
	Raw  string
	data any
	typ  reflect.Type
}

func (f *Field) UnmarshalJSON(data []byte) error {
	f.Raw = string(data)
	if f.Raw == "null" {
		f.Raw = ""
	}
	if strings.Contains(f.Raw, ".") {
		if n, err := strconv.ParseFloat(f.Raw, 64); err == nil {
			f.data = n
			f.typ = reflect.TypeOf(f.data)
		}
	} else if n, err := strconv.Atoi(string(data)); err == nil {
		f.data = n
		f.typ = reflect.TypeOf(n)
	} else if f.Raw == "true" || f.Raw == "false" {
		f.data = f.Raw == "true"
		f.typ = reflect.TypeOf(f.data)
	} else if strings.Contains(f.Raw, "e") {
		if n, err := strconv.ParseFloat(f.Raw, 64); err == nil {
			f.data = n
			f.typ = reflect.TypeOf(f.data)
		}
	}
	if f.data == nil {
		if f.Raw == "" {
			return nil
		}
		unquoted, err := strconv.Unquote(f.Raw)
		if err == nil {
			f.data = unquoted
			f.typ = reflect.TypeOf(f.data)
		}

	}
	return nil
}

func (f *Field) IsInt() bool {
	if f == nil || f.data == nil || f.typ == nil {
		return false
	}
	ok := f.typ.Kind() == reflect.Int
	if !ok && strings.HasPrefix(f.Raw, `"`) && strings.HasSuffix(f.Raw, `"`) {
		raw := strings.TrimPrefix(f.Raw, `"`)
		raw = strings.TrimSuffix(raw, `"`)
		_, err := strconv.Atoi(raw)
		ok = err == nil
	}
	return ok
}

func (f *Field) IsFloat64() bool {
	if f == nil || f.data == nil || f.typ == nil {
		return false
	}
	ok := f.typ.Kind() == reflect.Float64
	if !ok && strings.HasPrefix(f.Raw, `"`) && strings.HasSuffix(f.Raw, `"`) {
		raw := strings.TrimPrefix(f.Raw, `"`)
		raw = strings.TrimSuffix(raw, `"`)
		_, err := strconv.ParseFloat(raw, 64)
		ok = err == nil
	}
	return ok
}

func (f *Field) IsString() bool {
	if f == nil || f.data == nil || f.typ == nil {
		return false
	}
	return f.typ.Kind() == reflect.String
}

func (f *Field) Int() int {
	if f == nil || f.data == nil || f.typ == nil {
		return -1
	}
	ok := f.typ.Kind() == reflect.Int
	if ok {
		return f.data.(int)
	} else if strings.HasPrefix(f.Raw, `"`) && strings.HasSuffix(f.Raw, `"`) {
		raw := strings.TrimPrefix(f.Raw, `"`)
		raw = strings.TrimSuffix(raw, `"`)
		n, err := strconv.Atoi(raw)
		if err == nil {
			return n
		}
	}
	return -1
}

func (f *Field) BigInt() (*big.Int, bool) {
	if f == nil || f.data == nil || f.typ == nil {
		return nil, false
	}
	var bi big.Int
	raw := f.Raw
	if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
		raw = strings.TrimPrefix(raw, `"`)
		raw = strings.TrimSuffix(raw, `"`)
	}
	_, ok := bi.SetString(raw, 10)
	if !ok {
		return nil, false
	}
	return &bi, true
}

func (f *Field) Float64() float64 {
	if f == nil || f.data == nil || f.typ == nil {
		return -1
	}
	ok := f.typ.Kind() == reflect.Float64
	if ok {
		return f.data.(float64)
	} else if strings.HasPrefix(f.Raw, `"`) && strings.HasSuffix(f.Raw, `"`) {
		raw := strings.TrimPrefix(f.Raw, `"`)
		raw = strings.TrimSuffix(raw, `"`)
		n, err := strconv.ParseFloat(raw, 64)
		if err == nil {
			return n
		}
	}
	return -1
}

func (f *Field) BigFloat64() (*big.Float, bool) {
	if f == nil || f.data == nil || f.typ == nil {
		return nil, false
	}
	var bf big.Float
	raw := f.Raw
	if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
		raw = strings.TrimPrefix(raw, `"`)
		raw = strings.TrimSuffix(raw, `"`)
	}
	_, ok := bf.SetString(raw)
	if !ok {
		return nil, false
	}
	return &bf, true
}

func (f *Field) String() string {
	if f == nil || f.data == nil || f.typ == nil {
		return f.Raw
	}
	if f.IsString() && f.data != nil {
		return f.data.(string)
	}
	return f.Raw
}

func (f Field) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		fmt.Fprintf(s, "(%s) %s", f.typ, f.Raw)
	case 'q':
		fmt.Fprintf(s, "(%s) %q", f.typ, f.Raw)
	}
}

type StringField struct {
	Field
}
type IntField struct {
	Field
}

func (f *IntField) Int() int {
	if f.IsInt() {
		return f.Field.Int()
	}
	return -1
}

type Float64Field struct {
	Field
}

func (f *Float64Field) Float64() float64 {
	if f.IsFloat64() {
		return f.Field.Float64()
	}
	return -1
}
