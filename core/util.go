package core

type (
	form interface {
		Validate(Validator) error
	}
)

func validateForm(f form, v Validator) error {
	return f.Validate(v)
}

func boolPtr(v bool) *bool {
	return &v
}
