package validators

import (
	"fmt"
)

func (*intValidators) BetweenIncl(min, max int) func(i *int) error {
	return func(i *int) error {
		if err := isPresent(i); err != nil {
			return err
		}
		if *i < min || *i > max {
			return fmt.Errorf("value must be between %d and %d (incl.)", min, max)
		}
		return nil
	}
}

func (*intValidators) EqOrGreaterThan(min int) func(i *int) error {
	return func(i *int) error {
		if err := isPresent(i); err != nil {
			return err
		}
		if *i < min {
			return fmt.Errorf("value must be greater or equal than %d ", min)
		}
		return nil
	}
}
