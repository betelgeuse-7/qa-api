package validate

import "testing"

func TestValidate(t *testing.T) {
	v := New("email", "test@")
	v.Validate(
		IsEmail,
	)
}
