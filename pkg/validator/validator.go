package validator

import (
	"net/mail"
	"strings"
)

type Validator struct {
	Errors map[string]string `json:"errors"`
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Required(val, field string) {
	v.Check(strings.TrimSpace(val) != "", field, "this field is required")
}

func (v *Validator) MinLength(val string, min int, field string) {
	v.Check(len(strings.TrimSpace(val)) >= min, field, "must be at least minimum length")
}

func (v *Validator) Email(val, field string) {
	if strings.TrimSpace(val) == "" {
		v.AddError(field, "email cannot be empty")
		return
	}
	_, err := mail.ParseAddress(val)
	v.Check(err == nil, field, "invalid email format")
}

func (v *Validator) In(val string, allowed []string, field string) {
	for _, a := range allowed {
		if val == a {
			return
		}
	}
	v.AddError(field, "value must be one of allowed options")
}
