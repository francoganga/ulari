package ulari

import (
	"fmt"
	"testing"
)

func TestGenerateHTMLForm(t *testing.T) {
	data := struct {
		Name   string
		Age    int
		Active bool
	}{
		Name:   "John",
		Age:    30,
		Active: true,
	}

	form := generateHTMLForm(data)

	for _, v := range form.Values {

		switch v.(type) {
		case TextInput:
			fmt.Println("Text Input")
		case BoolInput:
			fmt.Println("Bool Input")
		case HiddenInput:
			fmt.Println("Hidden Input")
		case NumberInput:
			fmt.Println("Number Input")
		}
	}
}

func TestBind(t *testing.T) {
	data := struct {
		Name   string
		Age    int
		Active bool
	}{
		Name:   "John",
		Age:    30,
		Active: true,
	}

	form := generateHTMLForm(data)

	type Person struct {
		Name   string
		Age    int
		Active bool
	}

	person := Person{}

	err := form.Bind(&person)

	if err != nil {
		t.Fatal(err)
	}
}

