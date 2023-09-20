package validators

import (
	"errors"
	"fmt"
	"net/mail"
)

func (*stringValidators) NotEmpty(s *string) error {
	if err := isPresent(s); err != nil {
		return err
	}
	if *s == "" {
		return errors.New("empty value")
	}
	return nil
}

func (*stringValidators) ValidEmail(s *string) error {
	if err := isPresent(s); err != nil {
		return err
	}
	_, err := mail.ParseAddress(*s)
	return err
}

func (*stringValidators) Len(sl int) func(s *string) error {
	return func(s *string) error {
		if err := isPresent(s); err != nil {
			return err
		}
		if len(*s) != sl {
			return fmt.Errorf("string length must be %d", sl)
		}
		return nil
	}
}
