package ulari

import (
	"fmt"
	"reflect"
	"strconv"
)

type FormInput interface {
	GetValue() string
	GetName() string
}

type Form struct {
	Values []FormInput
}

type formInput struct {
	Name    string
	Value   string
	Classes []string
	Html    string
}

func (fi formInput) GetValue() string {
	return fi.Value
}

func (fi formInput) GetName() string {
	return fi.Name
}

type TextInput struct {
	formInput
}

type HiddenInput struct {
	formInput
}

type NumberInput struct {
	formInput
}

type BoolInput struct {
	formInput
}

func newForm() *Form {
	return &Form{
		Values: []FormInput{},
	}
}

func (f *Form) Bind(data interface{}) error {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	structVal := val.Elem()

	for _, v := range f.Values {
		_, ok := structVal.Type().FieldByName(v.GetName())

		if !ok {
			return fmt.Errorf("field %s not found in struct", v.GetName())
		}

		field := structVal.FieldByName(v.GetName())

		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("field %s is not valid or cannot be set", v.GetName())
		}

		switch t := v.(type) {
		case TextInput:
			field.SetString(t.GetValue())
		case HiddenInput:
		case BoolInput:
			field.SetBool(t.GetValue() == "on")
		case NumberInput:
			sv := t.GetValue()
			iv, err := strconv.ParseInt(sv, 10, 64)
			if err != nil {
				return fmt.Errorf("field %s is not a valid number: %w", sv, err)
			}
			field.SetInt(iv)
		}

	}

	return nil
}

func generateHTMLForm(data interface{}) *Form {

	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	form := newForm()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		switch field.Kind() {
		case reflect.String:
			fi := TextInput{
				formInput{
					Name:  fieldName,
					Value: field.String(),
				},
			}
			fi.Html = "<input type=\"text\" name=\"" + fieldName + "\" value=\"" + field.String() + "\">"
			form.Values = append(form.Values, fi)
		case reflect.Int:
			fi := NumberInput{
				formInput{
					Name:  fieldName,
					Value: field.String(),
				},
			}
			fi.Html = "<input type=\"number\" name=\"" + fieldName + "\" value=\"" + field.String() + "\">"
			form.Values = append(form.Values, fi)
		case reflect.Bool:
			fi := BoolInput{
				formInput{
					Name:  fieldName,
					Value: field.String(),
				},
			}
			s := "<input type=\"checkbox\" name=\"" + fieldName + "\""
			if field.Bool() {
				s += " checked>"
			}
			fi.Html = s
			form.Values = append(form.Values, fi)
		}
	}

	return form
}

