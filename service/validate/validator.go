package validate

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// a wrapper around validator.Validate. I want this to give me more useful error messages
type Validator struct {
	V *validator.Validate
	r *regexp.Regexp
}

type ValidationError string

const _USEFUL_ERROR_MESSAGE_REGEX = `key:\s.+:`

func New(v *validator.Validate) *Validator {
	r, err := regexp.Compile(_USEFUL_ERROR_MESSAGE_REGEX)
	if err != nil {
		panic("validate.New: regexp.Compile: " + err.Error())
	}
	return &Validator{V: v, r: r}
}

func (v *Validator) Validate(val interface{}) []ValidationError {
	res := []ValidationError{}
	err := v.V.Struct(val)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			errStr := strings.ToLower(e.Error())
			errStr = v.r.ReplaceAllString(errStr, "")
			res = append(res, ValidationError(errStr))
		}
	}
	return res
}
