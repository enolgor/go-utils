package examples

import (
	"errors"
	"os"
	"time"

	"github.com/enolgor/go-utils/conf"
	"github.com/enolgor/go-utils/validators"
	"golang.org/x/text/language"
)

var PORT int
var HOST string
var TIMEZONE time.Location
var TEST bool

var LANG conf.KeyValue[language.Tag, bool]

func init() {
	os.Setenv("HOST", "asdf")
	conf.SetValidate(&PORT, "PORT", "p", 8080, validators.Integers.EqOrGreaterThan(10000))
	conf.SetEnv(&HOST, "HOST", "localhost")
	conf.SetFlag(&TIMEZONE, "tz", *time.UTC)
	conf.Set(&TEST, "TEST", "t", false)
	conf.SetPairValidate(&LANG, "LANGUAGE", "lang", conf.KeyValue[language.Tag, bool]{Key: language.English, Value: true}, keyValidator, valueValidator, keyValueValidator)
	conf.Read()
}

func keyValidator(key *language.Tag) error {
	if key == nil {
		return errors.New("nil key")
	}
	if key.String() == "es" {
		return errors.New("spanish not allowed")
	}
	return nil
}

func valueValidator(value *bool) error {
	if value == nil {
		return errors.New("nil value")
	}
	if !*value {
		return errors.New("value can only be true")
	}
	return nil
}

func keyValueValidator(key *language.Tag, value *bool) error {
	if key == nil || value == nil {
		return errors.New("nil key or value")
	}
	if key.String() == "fr" && *value {
		return errors.New("french can only be false")
	}
	return nil
}
