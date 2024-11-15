package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Maximum length for this field is %d", d))
	}
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("Minimum length for this field is %d", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) EmailValid(field string) {
	_, err := mail.ParseAddress(f.Get("email"))
	if err != nil {
		f.Errors.Add(field, "Email address is not valid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) GetInstance(dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a non-nil pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("dst must be a pointer to a struct")
	}

	for fieldName, value := range f.Values {
		field := v.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, fieldName) })
		if field.IsValid() && field.CanSet() {
			field.SetString(value[0])
		}
	}

	return nil
}
