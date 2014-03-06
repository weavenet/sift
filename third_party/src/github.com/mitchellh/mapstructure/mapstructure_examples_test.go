package mapstructure

import (
	"fmt"
)

func ExampleDecode() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var result Person
	err := Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result)
	// Output:
	// mapstructure.Person{Name:"Mitchell", Age:91, Emails:[]string{"one", "two", "three"}, Extra:map[string]string{"twitter":"mitchellh"}}
}

func ExampleDecode_errors() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   123,
		"age":    "bad value",
		"emails": []int{1, 2, 3},
	}

	var result Person
	err := Decode(input, &result)
	if err == nil {
		panic("should have an error")
	}

	fmt.Println(err.Error())
	// Output:
	// 5 error(s) decoding:
	//
	// * 'Name' expected type 'string', got unconvertible type 'int'
	// * 'Age' expected type 'int', got unconvertible type 'string'
	// * 'Emails[0]' expected type 'string', got unconvertible type 'int'
	// * 'Emails[1]' expected type 'string', got unconvertible type 'int'
	// * 'Emails[2]' expected type 'string', got unconvertible type 'int'
}

func ExampleDecode_metadata() {
	type Person struct {
		Name string
		Age  int
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":  "Mitchell",
		"age":   91,
		"email": "foo@bar.com",
	}

	// For metadata, we make a more advanced DecoderConfig so we can
	// more finely configure the decoder that is used. In this case, we
	// just tell the decoder we want to track metadata.
	var md Metadata
	var result Person
	config := &DecoderConfig{
		Metadata: &md,
		Result:   &result,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		panic(err)
	}

	if err := decoder.Decode(input); err != nil {
		panic(err)
	}

	fmt.Printf("Unused keys: %#v", md.Unused)
	// Output:
	// Unused keys: []string{"email"}
}

func ExampleDecode_weaklyTypedInput() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON, generated by a weakly typed language
	// such as PHP.
	input := map[string]interface{}{
		"name":   123,                      // number => string
		"age":    "42",                     // string => number
		"emails": map[string]interface{}{}, // empty map => empty array
	}

	var result Person
	config := &DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &result,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result)
	// Output: mapstructure.Person{Name:"123", Age:42, Emails:[]string{}}
}
