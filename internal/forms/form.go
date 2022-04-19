package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// New returns a new Form struct with given data.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks whether every required field is not empty.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

func (f *Form) MinLength(field string, length int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters).", length))
	}
}

func (f *Form) MaxLength(field string, length int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > length {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters).", length))
	}
}

func (f *Form) MatchPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid.")
	}
}

// IsValid returns true if there is are no errors.
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
