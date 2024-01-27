package ulari

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type FormInput interface {
	GetValue() string
	GetName() string
}

type Form struct {
	Fields []FormInput
	Data   *url.Values
}

func (f *Form) NameValid(name string) bool {
	for _, v := range f.Fields {
		if v.GetName() == name {
			return true
		}
	}

	return false
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
		Fields: []FormInput{},
	}
}

func newFromData(data *url.Values) *Form {

	form := newForm()

	form.Data = data

	return form
}

func (f *Form) Bind(data interface{}) error {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	structVal := val.Elem()

	for k, v := range *f.Data {
		fmt.Printf("v=%v\n", v)
		_, ok := structVal.Type().FieldByName(k)

		if !ok {
			return fmt.Errorf("field %s not found in struct", k)
		}

		field := structVal.FieldByName(k)

		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("field %s is not valid or cannot be set", k)
		}

		switch field.Kind() {

		case reflect.String:
			field.SetString(v[0])

		case reflect.Int:
			iv, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return fmt.Errorf("field %s is not a valid number: %w", v[0], err)
			}
			field.SetInt(iv)
		case reflect.Bool:
			field.SetBool(v[0] == "on")
		default:
			return fmt.Errorf("field %s is not a valid type", k)
		}

		// switch t := v.(type) {
		// case TextInput:
		// 	field.SetString(t.GetValue())
		// case HiddenInput:
		// case BoolInput:
		// 	field.SetBool(t.GetValue() == "on")
		// case NumberInput:
		// 	sv := t.GetValue()
		// 	iv, err := strconv.ParseInt(sv, 10, 64)
		// 	if err != nil {
		// 		return fmt.Errorf("field %s is not a valid number: %w", sv, err)
		// 	}
		// 	field.SetInt(iv)
		// }

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
			fi.Html = fmt.Sprintf(`<input type="text" name="%s" value="%s">`, fieldName, field.String())

			form.Fields = append(form.Fields, fi)
		case reflect.Int:
			fi := NumberInput{
				formInput{
					Name:  fieldName,
					Value: field.String(),
				},
			}
			fi.Html = fmt.Sprintf(`<input type="number" name="%s" value="%s">`, fieldName, field.String())
			form.Fields = append(form.Fields, fi)
		case reflect.Bool:
			fi := BoolInput{
				formInput{
					Name:  fieldName,
					Value: field.String(),
				},
			}
			s := fmt.Sprintf(`<input type="checkbox" name="%s" value="%s">`, fieldName, field.String())
			if field.Bool() {
				s += " checked>"
			}
			fi.Html = s
			form.Fields = append(form.Fields, fi)
		}
	}

	return form
}

