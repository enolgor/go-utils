package validators

import (
	"errors"

	"github.com/enolgor/go-utils/parse"
)

type stringValidators struct{}
type intValidators struct{}
type hexValidators struct{}

var Strings *stringValidators = &stringValidators{}
var Ints *intValidators = &intValidators{}
var HexBytes *hexValidators = &hexValidators{}

func isPresent[P parse.Parseable](p *P) error {
	if p == nil {
		return errors.New("nil value")
	}
	return nil
}

func All[P parse.Parseable](validators ...func(*P) error) func(*P) error {
	return func(v *P) error {
		var errs error
		for i := range validators {
			errs = errors.Join(errs, validators[i](v))
		}
		return errs
	}
}

func Any[P parse.Parseable](validators ...func(*P) error) func(*P) error {
	return func(v *P) error {
		var errs error
		for i := range validators {
			err := validators[i](v)
			if err != nil {
				return nil
			}
			errs = errors.Join(errs)
		}
		return errs
	}
}
