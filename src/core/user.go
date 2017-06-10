package core

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"sync"

	"github.com/olegakbarov/io.confs.core/src/domain"
	"github.com/pkg/errors"
)

type (
	TokenType uint8

	// User service interface
	User interface {
		GetFromAuthToken(tokenStr string) (*domain.User, error)
		GenToken(*domain.User, TokenType) (string, error)
		Login(*LoginRequest) (*domain.User, error)
		Register(*RegisterRequest) (*domain.User, error)
		Activate(*UserActivateRequest) error
		SendPasswordResetMail(*ForgotPasswordRequest) error
		ResetPassword(*ResetPasswordRequest) error
		Show(*ShowUserRequest) (*domain.User, error)
		Update(*UpdateUserRequest) error
	}

	user struct {
		jwt       JWTSignParser
		repo      UserRepository
		mailer    Mailer
		validator Validator
		emitter   Emitter
	}

	// LoginRequest context for user.Login()
	LoginRequest struct {
		Email    string
		Password string
	}

	// RegisterRequest context for user.Register()
	RegisterRequest struct {
		Email         string
		Password      string
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		IsActive      *bool  `json:"-"`
		ActivationURL string `json:"-"`
	}

	// UserActivateRequest context for user.Activate()
	UserActivateRequest struct {
		Token string `json:"token"`
	}

	//ForgotPasswordRequest context for user.ForgotPassword()
	ForgotPasswordRequest struct {
		Email   string
		BaseURL string
	}

	// ResetPasswordRequest context for user.ResetPassword()
	ResetPasswordRequest struct {
		Token    string
		Password string
	}

	// ShowUserRequest context for user.Show()
	ShowUserRequest struct {
		ID uint
	}

	// UpdateUserRequest context for user.Update()
	UpdateUserRequest struct {
		ID uint `json:"-"`
		RegisterRequest
	}
)

var (
	userInstance User
	userOnce     sync.Once
)

// Token Types
const (
	AuthToken TokenType = iota
	ActivationToken
	PasswordResetToken
)

func (f *factory) NewUser() User {
	userOnce.Do(func() {
		userInstance = &user{
			jwt:       f.jwt,
			repo:      f.NewUserRepository(),
			mailer:    f.NewMail(),
			validator: f.v,
			emitter:   f.emitter,
		}
	})

	return userInstance
}

func (u *user) Login(r *LoginRequest) (*domain.User, error) {
	usr, err := u.repo.OneByEmail(r.Email)
	if err != nil {
		if err == ErrNoRows {
			return nil, ErrWrongCredentials
		}
		return nil, err
	}

	if !usr.IsCredentialsVerified(r.Password) {
		return nil, ErrWrongCredentials
	}

	if !*usr.IsActive {
		return nil, ErrInActiveUser
	}
	return usr, nil
}

func (u *user) Register(r *RegisterRequest) (*domain.User, error) {
	// validation
	if err := u.validator.CheckEmail(r.Email); err != nil {
		return nil, err
	}
	if err := checkPassword(u.validator, r.Password); err != nil {
		return nil, err
	}

	// check for email
	exists, err := u.repo.ExistsByEmail(r.Email)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, ErrEmailExists
	}

	if r.IsActive == nil {
		r.IsActive = boolPtr(false)
	}

	var usr domain.User
	usr.FirstName = r.FirstName
	usr.LastName = r.LastName
	usr.Email = r.Email
	usr.SetPassword(r.Password)
	usr.IsActive = r.IsActive
	usr.IsAdmin = boolPtr(false)

	if err := u.repo.Add(&usr); err != nil {
		return nil, err
	}

	// so we need to send activation mail?
	if *r.IsActive {
		return &usr, nil // no :)
	}

	token, err := u.GenToken(&usr, ActivationToken)
	if err != nil {
		return nil, err
	}

	u.emitter.Emit(TokenGenerated, token, ActivationToken)

	if r.ActivationURL == "" {
		r.ActivationURL = os.Getenv("ACTIVATION_URL")
	}
	pURL, err := url.Parse(r.ActivationURL)
	if err != nil {
		return nil, err
	}
	q := pURL.Query()
	q.Set("token", token)
	pURL.RawQuery = q.Encode()

	if err := u.mailer.SendWelcomeMail(usr.Email, pURL.String()); err != nil {
		return nil, err
	}

	return &usr, nil
}

func (u *user) Activate(r *UserActivateRequest) error {
	email, err := u.getEmailFromToken(r.Token, ActivationToken)
	if err != nil {
		return err
	}

	usr, err := u.repo.OneByEmail(email)
	if err != nil {
		return err
	}

	// check if already active?
	if *usr.IsActive {
		return nil
	}

	// activate user
	usr.IsActive = boolPtr(true)

	return u.repo.Update(usr)
}

func (u *user) GetFromAuthToken(tokenStr string) (*domain.User, error) {
	email, err := u.getEmailFromToken(tokenStr, AuthToken)
	if err != nil {
		return nil, err
	}

	return u.repo.OneByEmail(email)
}

func (u *user) SendPasswordResetMail(r *ForgotPasswordRequest) error {
	if r.BaseURL == "" {
		r.BaseURL = os.Getenv("PASSWORD_RESET_URL")
	}

	usr, err := u.repo.OneByEmail(r.Email)
	if err != nil {
		return err
	}

	token, err := u.GenToken(usr, PasswordResetToken)
	if err != nil {
		return err
	}

	pURL, err := url.Parse(r.BaseURL)
	if err != nil {
		return err
	}
	q := pURL.Query()
	q.Set("token", token)
	pURL.RawQuery = q.Encode()

	return u.mailer.SendPasswordResetLink(r.Email, pURL.String())
}

func (u *user) ResetPassword(r *ResetPasswordRequest) error {
	// validation
	if err := checkPassword(u.validator, r.Password); err != nil {
		return err
	}

	email, err := u.getEmailFromToken(r.Token, PasswordResetToken)
	if err != nil {
		return err
	}

	usr, err := u.repo.OneByEmail(email)
	if err != nil {
		return err
	}

	usr.SetPassword(r.Password)

	return u.repo.Update(usr)
}

func (u *user) Show(r *ShowUserRequest) (*domain.User, error) {
	return u.repo.One(r.ID)
}

func (u *user) Update(r *UpdateUserRequest) error {
	usr, err := u.repo.One(r.ID)
	if err != nil {
		return err
	}

	if r.FirstName != "" {
		usr.FirstName = r.FirstName
	}
	if r.LastName != "" {
		usr.LastName = r.LastName
	}
	if r.Email != "" && r.Email != usr.Email {
		// validate
		if err := u.validator.CheckEmail(r.Email); err != nil {
			return err
		}
		// someone else exists with this email?
		exists, err := u.repo.ExistsByEmail(r.Email)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailExists
		}
		usr.Email = r.Email
	}
	if r.Password != "" {
		if err := checkPassword(u.validator, r.Password); err != nil {
			return err
		}
		usr.SetPassword(r.Password)
	}

	return u.repo.Update(usr)
}

func (u *user) GenToken(usr *domain.User, t TokenType) (string, error) {
	claims := map[string]interface{}{
		"type":  t,
		"email": usr.Email,
	}
	switch t {
	case AuthToken:
		claims["exp"] = time.Now().Add(time.Hour * 6).Unix()
	case ActivationToken:
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	case PasswordResetToken:
		claims["exp"] = time.Now().Add(time.Hour * 3).Unix()
	default:
		return "", errors.Errorf("undefined token type %v", t)
	}
	return u.jwt.Sign(claims, os.Getenv("SECRET_KEY"))
}

func (u *user) getEmailFromToken(token string, t TokenType) (string, error) {
	claims, err := u.jwt.Parse(token, os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}

	if ct, ok := claims["type"].(float64); ok != true || TokenType(ct) != t {
		return "", &TokenErr{fmt.Sprintf("invalid token type %v", t), false}
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("email can't get from token claims: %v", claims)
	}

	return email, nil
}

func checkPassword(v Validator, p string) error {
	return v.CheckStringLen(p, 4, 16, "Password")
}
