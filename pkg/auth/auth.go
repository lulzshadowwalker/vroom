package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/lulzshadowwalker/vroom/internal/database/model"
	"github.com/manifoldco/promptui"
)

func init() {
	validate = validator.New()
}

var validate *validator.Validate

type AuthHandler struct {
	Repo AuthHandlerRepo
}

type AuthHandlerRepo interface {
	SignIn(email, password string) (*model.User, error)
	SignUp(username, email, password string) (*model.User, error)
}

func (a *AuthHandler) Trigger() (*model.User, error) {
	type opt string

	const (
		signInOpt opt = "Sign in"
		signUpOpt opt = "Sign up"
		guestOpt  opt = "Continue as a Guest"
	)

	prompt := promptui.Select{
		Label: "✨ Authentication",
		Items: []opt{signInOpt, signUpOpt, guestOpt},
	}

	_, res, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("cannot run authentication prompt %w", err)
	}

	switch opt(res) {
	case signInOpt:
		return a.signIn()
	case signUpOpt:
		return a.signUp()
	case guestOpt:
		panic("unimplemented :: guest auth")
	default:
		return nil, errors.New("invalid authentication option")
	}
}

func (a *AuthHandler) signIn() (*model.User, error) {
	email, err := getEmail()
	if err != nil {
		return nil, err
	}

	password, err := getPassword()
	if err != nil {
		return nil, err
	}

	user, err := a.Repo.SignIn(email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthHandler) signUp() (*model.User, error) {
	uname, err := getUsername()
	if err != nil {
		return nil, err
	}

	email, err := getEmail()
	if err != nil {
		return nil, err
	}

	pwd, err := getPassword()
	if err != nil {
		return nil, err
	}

	user, err := a.Repo.SignUp(uname, email, pwd)
	if err != nil {
		return nil, fmt.Errorf("cannot sign up %w", err)
	}

	return user, nil
}

func getEmail() (res string, err error) {
	validate := func(input string) error {
		err := validate.Var(input, "required,email")
		if err != nil {
			return errors.New("where is the email pepega")
		}

		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "email",
		Templates: templates,
		Validate:  validate,
	}

	res, err = prompt.Run()
	if err != nil {
		return "", fmt.Errorf("cannot run email prompt %w", err)
	}

	return strings.Trim(res, " "), nil
}

func getPassword() (res string, err error) {
	validate := func(input string) error {
		err := validate.Var(input, "required,min=8")
		if err != nil {
			return errors.New("where is the password pepega")
		}

		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "password (min 8 chars)",
		Templates: templates,
		Validate:  validate,
		Mask:      '⏺',
	}

	res, err = prompt.Run()
	if err != nil {
		return "", fmt.Errorf("cannot run password prompt %w", err)
	}

	return strings.Trim(res, " "), nil
}

func getUsername() (res string, err error) {
	validate := func(input string) error {
		err := validate.Var(input, "required,min=3")
		if err != nil {
			return errors.New("where is the username pepega")
		}

		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "username (min 3 chars)",
		Templates: templates,
		Validate:  validate,
	}

	res, err = prompt.Run()
	if err != nil {
		return "", fmt.Errorf("cannot run username prompt %w", err)
	}

	return strings.Trim(res, " "), nil
}
