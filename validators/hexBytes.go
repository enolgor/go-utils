package validators

import (
	"fmt"

	"github.com/enolgor/go-utils/parse/types"
)

func (*hexValidators) NotEmpty(v *types.HexBytes) error {
	if err := isPresent(v); err != nil {
		return err
	}
	if len(*v) == 0 {
		return fmt.Errorf("empty byte array")
	}
	return nil
}

func (*hexValidators) Len(sl int) func(v *types.HexBytes) error {
	return func(v *types.HexBytes) error {
		if err := isPresent(v); err != nil {
			return err
		}
		if len(*v) != sl {
			return fmt.Errorf("bytes length must be %d", sl)
		}
		return nil
	}
}
