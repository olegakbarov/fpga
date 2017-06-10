package providers

// "github.com/alioygur/gocart/engine"
// "github.com/alioygur/is"

// var newErr = engine.NewValidationErr

// func NewValidator() engine.Validator {
//     return &validator{}
// }

type validator struct{}

// CheckEmail checks input whatever email or not
// func (v *validator) CheckEmail(email string) error {
//     _, err := mail.ParseAddress(email)
//     if err != nil {
//         return newErr("invalid email address")
//     }
//     return nil
// }

// func (v *validator) CheckRequired(val, field string) error {
//     if len(val) != 0 {
//         return nil
//     }
//     return newErr("the %s field is required", field)
// }

// func (v *validator) CheckStringLen(val string, min int, max int, field string) error {
//     if is.StringLength(val, min, max) {
//         return nil
//     }
//     return newErr("the %s field length must between %d and %d", field, min, max)
// }
