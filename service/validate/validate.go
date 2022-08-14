// After being unsuccessful to find a good validator package, I decided to roll my own little one.
package validate

import (
	"fmt"
	"strconv"
	"strings"
)

type opts struct {
	isEmail, hasLengthBounds, hasMinCheck, hasMaxCheck    bool
	lengthBounds                                          []int
	minNumVal, maxNumVal                                  uint
	doesNotStartWith, startsWith, doesNotExist, doesExist rune
}

type Validator struct {
	valName string
	value   string
	opts    opts
}

// name == field name (e.g. handle, username, password ...)
func New(name string, value string) *Validator {
	v := &Validator{value: value, valName: name}
	return v
}

type OptFunc func(*Validator)

func (v *Validator) Validate(opts ...OptFunc) []string {
	res := []string{}
	for _, f := range opts {
		f(v)
	}
	if v.opts.doesExist != rune(0) {
		res = append(res, v.mustContain(v.opts.doesExist))
	}
	if v.opts.doesNotExist != rune(0) {
		res = append(res, v.mustNotContain(v.opts.doesNotExist))
	}
	if v.opts.doesNotStartWith != rune(0) {
		res = append(res, v.doesNotStartWith(v.opts.doesNotStartWith))
	}
	if v.opts.startsWith != rune(0) {
		res = append(res, v.doesStartWith(v.opts.startsWith))
	}
	if v.opts.isEmail {
		res = append(res, v.isEmail())
	}
	if v.opts.hasLengthBounds {
		res = append(res, v.assertLength(v.opts.lengthBounds[0], v.opts.lengthBounds[1]))
	}
	if v.opts.hasMinCheck {
		res = append(res, v.min(v.opts.minNumVal))
	}
	if v.opts.hasMaxCheck {
		res = append(res, v.max(v.opts.maxNumVal))
	}
	return res
}

type validationError = string

func (v *Validator) isEmail() validationError {
	err_ := validationError(v.valName + ": invalid email")
	val := v.value
	if val[0] == '@' {
		return err_
	}
	if !(strings.Contains(val, "@")) {
		return err_
	}
	if !(strings.Contains(val, ".")) {
		return err_
	}
	return validationError("")
}

func (v *Validator) assertLength(min, max int) validationError {
	val := v.value
	name := v.valName
	length := len(val)
	if length < min {
		return validationError(fmt.Sprintf("%s: minimum length is %d", name, min))
	}
	if max != -1 {
		if length > max {
			return validationError(fmt.Sprintf("%s: maximum length is %d", name, max))
		}
	}
	return ""
}

func (v *Validator) min(n uint) validationError {
	val := v.value
	name := v.valName
	valUint64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return validationError(name + ": invalid input")
	}
	if uint(valUint64) < n {
		return validationError(fmt.Sprintf("%s: less than minimum value of %d", name, n))
	}
	return ""
}

func (v *Validator) max(n uint) validationError {
	val := v.value
	name := v.valName
	valUint64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return validationError(name + ": invalid input")
	}
	if uint(valUint64) > n {
		return validationError(fmt.Sprintf("%s: bigger than maximum value of %d", name, n))
	}
	return ""
}

func (v *Validator) doesStartWith(char rune) validationError {
	if []rune(v.value)[0] != char {
		return validationError(fmt.Sprintf("%s: first character is not '%v'", v.valName, char))
	}
	return ""
}

func (v *Validator) doesNotStartWith(char rune) validationError {
	err := v.doesStartWith(char)
	if err == "" {
		return validationError(fmt.Sprintf("%s: first character must not be '%v'", v.valName, char))
	}
	return ""
}

func (v *Validator) mustContain(char rune) validationError {
	if ok := strings.Contains(v.value, string(char)); !(ok) {
		return validationError(fmt.Sprintf("%s: must contain '%v'", v.valName, char))
	}
	return ""
}

func (v *Validator) mustNotContain(char rune) validationError {
	err := v.mustContain(char)
	if err == "" {
		return validationError(fmt.Sprintf("%s: must not contain '%v'", v.valName, char))
	}
	return ""
}

func IsEmail(v *Validator) {
	v.opts.isEmail = true
}

func StringLength(v *Validator, min, max int) {
	v.opts.hasLengthBounds = true
	v.opts.lengthBounds = []int{min, max}
}

func Min(v *Validator, n uint) {
	v.opts.hasMinCheck = true
	v.opts.minNumVal = n
}

func Max(v *Validator, n uint) {
	v.opts.hasMaxCheck = true
	v.opts.maxNumVal = n
}

func DoesNotStartWith(v *Validator, char rune) {
	v.opts.doesNotStartWith = char
}

func MustStartWith(v *Validator, char rune) {
	v.opts.startsWith = char
}

func MustNotContain(v *Validator, char rune) {
	v.opts.doesNotExist = char
}

func MustContain(v *Validator, char rune) {
	v.opts.doesExist = char
}
