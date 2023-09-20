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
