package conf

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/enolgor/go-utils/parse"
)

type KeyValue[K Configurable, V Configurable] struct {
	Key   K
	Value V
}

type Configurable interface {
	parse.Parseable
}

func identity(s string) (string, error) {
	return s, nil
}

func flagFunc[T Configurable](take *T, f func(string) (string, error)) func(string) error {
	if f == nil {
		f = identity
	}
	return func(s string) error {
		var err error
		if s, err = f(s); err != nil {
			return err
		}
		return parse.Parse(take, s)
	}
}

func chain(f1 func(string) error, f2 func(string) error) func(string) error {
	return func(s string) error {
		var err error
		if err = f1(s); err != nil {
			return err
		}
		return f2(s)
	}
}

func getEnv[T Configurable](take *T, key string, f func(string) (string, error)) (string, error) {
	var err error
	var str string
	var ok bool
	if f == nil {
		f = identity
	}
	if str, ok = os.LookupEnv(key); ok && str != "" {
		if str, err = f(str); err != nil {
			return str, err
		}
		err = parse.Parse(take, str)
	}
	return str, err
}

var envErr error

var flagFuncs map[string]func(string) error = map[string]func(string) error{}
var boolFlags map[string]*bool = map[string]*bool{}

func set[T Configurable](take *T, envKey, flagKey string, def T, f func(string) (string, error)) {
	*take = def
	if envKey != "" {
		if val, err := getEnv(take, envKey, f); err != nil {
			envErr = fmt.Errorf(`invalid value "%s" for env "%s": %s`, val, envKey, err.Error())
		}
	}
	if flagKey != "" {
		if v, ok := any(take).(*bool); ok && f == nil {
			boolFlags[flagKey] = v
		} else {
			if fn, ok := flagFuncs[flagKey]; ok {
				flagFuncs[flagKey] = chain(fn, flagFunc(take, f))
			} else {
				flagFuncs[flagKey] = flagFunc(take, f)
			}
		}
	}
}

func setKVarray[K Configurable, V Configurable](take *[]*KeyValue[K, V], envKey, flagKey string, def []KeyValue[K, V], i int) {
	//fmt.Printf("called i=%d\n", i)
	if len(*take) < i+1 {
		*take = append(*take, &KeyValue[K, V]{})
	}
	set[K](&((*take)[i].Key), envKey, flagKey, def[0].Key, func(str string) (string, error) {
		parts := strings.Split(str, ",")
		//fmt.Printf("set=%d,len=%d\n", i, len(parts))
		if len(parts) != i+1 {
			//fmt.Println("calling kv array")
			setKVarray(take, envKey, flagKey, def, i+1)
		}
		str = parts[i]
		//fmt.Println(str)
		idx := strings.Index(str, "=")
		if idx == -1 {
			return "", fmt.Errorf(`"=" sign not found for key-value property`)
		}
		return str[:idx], nil
	})
	set[V](&((*take)[i].Value), envKey, flagKey, def[0].Value, func(str string) (string, error) {
		parts := strings.Split(str, ",")
		str = parts[i]
		//fmt.Println(str)
		idx := strings.Index(str, "=")
		if idx == -1 {
			return "", fmt.Errorf(`"=" sign not found for key-value property`)
		}
		return str[idx+1:], nil
	})
}

// func SetKVs[K Configurable, V Configurable](take *[]*KeyValue[K, V], envKey, flagKey string, def []KeyValue[K, V]) {
// 	setKVarray(take, envKey, flagKey, def, 0)
// }

var validators []func() error = []func() error{}

func wrap[T Configurable](env, flag string, take *T, validator func(*T) error) {
	validators = append(validators, func() error {
		if err := validator(take); err != nil {
			return fmt.Errorf(`invalid value for env "%s" or flag "%s", %s`, env, flag, err.Error())
		}
		return nil
	})
}

func wrapKV[K Configurable, V Configurable](env, flag string, take *KeyValue[K, V], keyValidator func(*K) error, valueValidator func(*V) error, keyValueValidator func(*K, *V) error) {
	validators = append(validators, func() error {
		if err := keyValidator(&take.Key); err != nil {
			return fmt.Errorf(`invalid key for env "%s" or flag "%s", %s`, env, flag, err.Error())
		}
		return nil
	})
	validators = append(validators, func() error {
		if err := valueValidator(&take.Value); err != nil {
			return fmt.Errorf(`invalid value for env "%s" or flag "%s", %s`, env, flag, err.Error())
		}
		return nil
	})
	validators = append(validators, func() error {
		if err := keyValueValidator(&take.Key, &take.Value); err != nil {
			return fmt.Errorf(`invalid keyValue for env "%s" or flag "%s", %s`, env, flag, err.Error())
		}
		return nil
	})
}

func SetValidate[T Configurable](take *T, envKey, flagKey string, def T, validator func(*T) error) {
	wrap(envKey, flagKey, take, validator)
	set[T](take, envKey, flagKey, def, nil)
}

func SetPairValidate[K Configurable, V Configurable](take *KeyValue[K, V], envKey, flagKey string, def KeyValue[K, V], keyValidator func(*K) error, valueValidator func(*V) error, keyValueValidator func(*K, *V) error) {
	wrapKV(envKey, flagKey, take, keyValidator, valueValidator, keyValueValidator)
	*take = def
	setKVarray(&[]*KeyValue[K, V]{take}, envKey, flagKey, []KeyValue[K, V]{*take}, 0)
}

func SetEnvValidate[T Configurable](take *T, envKey string, def T, validator func(*T) error) {
	SetValidate(take, envKey, "", def, validator)
}

func SetFlagValidate[T Configurable](take *T, flagKey string, def T, validator func(*T) error) {
	SetValidate(take, "", flagKey, def, validator)
}

func nopValidate[T Configurable](t *T) error {
	if t == nil {
		return errors.New("value is nill")
	}
	return nil
}

func nopKVValidate[K Configurable, V Configurable](k *K, v *V) error {
	if k == nil {
		return errors.New("key is nill")
	}
	if v == nil {
		return errors.New("value is nill")
	}
	return nil
}

func Set[T Configurable](take *T, envKey, flagKey string, def T) {
	SetValidate(take, envKey, flagKey, def, nopValidate)
}

func SetPair[K Configurable, V Configurable](take *KeyValue[K, V], envKey, flagKey string, def KeyValue[K, V]) {
	SetPairValidate(take, envKey, flagKey, def, nopValidate, nopValidate, nopKVValidate)
}

func SetEnv[T Configurable](take *T, envKey string, def T) {
	SetEnvValidate(take, envKey, def, nopValidate)
}

func SetFlag[T Configurable](take *T, flagKey string, def T) {
	SetFlagValidate(take, flagKey, def, nopValidate)
}

// func SetEnvKV[K Configurable, V Configurable](take *KeyValue[K, V], envKey string, def KeyValue[K, V]) {
// 	SetKV(take, envKey, "", def)
// }

// func SetFlagKV[K Configurable, V Configurable](take *KeyValue[K, V], flagKey string, def KeyValue[K, V]) {
// 	SetKV(take, "", flagKey, def)
// }

func Read() {
	if envErr != nil {
		panic(envErr)
	}
	for key, v := range boolFlags {
		if _, ok := flagFuncs[key]; !ok {
			flag.BoolVar(v, key, *v, "")
		}
	}
	for key, fn := range flagFuncs {
		flag.Func(key, "", fn)
	}
	boolFlags = map[string]*bool{}
	flagFuncs = map[string]func(string) error{}
	flag.Parse()
	for _, validator := range validators {
		if err := validator(); err != nil {
			panic(err)
		}
	}
	validators = []func() error{}
}
