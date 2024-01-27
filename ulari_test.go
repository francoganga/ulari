package ulari

import (
	"fmt"
	"net/url"
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

	fmt.Printf("form=%+v\n", form)

	for _, v := range form.Fields {

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

	data := make(url.Values)
	data["Name"] = []string{"John"}
	data["Age"] = []string{"30"}
	data["Active"] = []string{"true"}

	form := newFromData(&data)

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

	fmt.Printf("person=%+v\n", person)
}

