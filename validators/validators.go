package validators

import (
	"errors"
	"fmt"
	"net/mail"
)

type validator struct {
	Strings  stringValidators
	Integers intValidators
}

type stringValidators struct {
	NotEmpty   func(s *string) error
	ValidEmail func(s *string) error
	Len        func(sl int) func(s *string) error
}

type intValidators struct {
	BetweenIncl     func(min, max int) func(i *int) error
	EqOrGreaterThan func(min int) func(i *int) error
}

var Validator validator = validator{
	Strings: stringValidators{
		NotEmpty:   stringNotEmpty,
		ValidEmail: stringValidEmail,
		Len:        stringLen,
	},
	Integers: intValidators{
		BetweenIncl:     intBetweenIncl,
		EqOrGreaterThan: intEqOrGreaterThan,
	},
}

func stringNotEmpty(s *string) error {
	if s == nil {
		return errors.New("nil value")
	}
	if *s == "" {
		return errors.New("empty value")
	}
	return nil
}

func stringValidEmail(s *string) error {
	if s == nil {
		return errors.New("nil email")
	}
	_, err := mail.ParseAddress(*s)
	return err
}

func stringLen(sl int) func(s *string) error {
	return func(s *string) error {
		if s == nil {
			return errors.New("nil value")
		}
		if len(*s) != sl {
			return fmt.Errorf("string length must be %d", sl)
		}
		return nil
	}
}

func intBetweenIncl(min, max int) func(i *int) error {
	return func(i *int) error {
		if i == nil {
			return errors.New("nil value")
		}
		if *i < min || *i > max {
			return fmt.Errorf("value must be between %d and %d (incl.)", min, max)
		}
		return nil
	}
}

func intEqOrGreaterThan(min int) func(i *int) error {
	return func(i *int) error {
		if i == nil {
			return errors.New("nil value")
		}
		if *i < min {
			return fmt.Errorf("value must be greater or equal than %d ", min)
		}
		return nil
	}
}
