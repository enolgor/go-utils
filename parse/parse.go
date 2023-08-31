package parse

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/enolgor/go-utils/parse/types"
	"golang.org/x/text/language"
)

type Parseable interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		bool | string |
		complex64 | complex128 |
		time.Duration | time.Time | time.Location | language.Tag |
		[]int | []int8 | []int16 | []int32 | []int64 |
		[]uint | []uint8 | []uint16 | []uint32 | []uint64 |
		[]float32 | []float64 |
		[]bool | []string |
		[]complex64 | []complex128 |
		[]time.Duration | []time.Time | []time.Location | []language.Tag |
		types.HexByte | types.OctByte |
		types.HexBytes | types.B32Bytes | types.B64Bytes |
		[]types.HexByte | []types.OctByte |
		[]types.HexBytes | []types.B32Bytes | []types.B64Bytes
}

func Int(str string) (int, error) {
	return strconv.Atoi(str)
}

func IntToString(v int) string {
	return strconv.Itoa(v)
}

func Int8(str string) (int8, error) {
	v, err := strconv.ParseInt(str, 10, 8)
	return int8(v), err
}

func Int8ToString(v int8) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int16(str string) (int16, error) {
	v, err := strconv.ParseInt(str, 10, 16)
	return int16(v), err
}

func Int16ToString(v int16) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int32(str string) (int32, error) {
	v, err := strconv.ParseInt(str, 10, 32)
	return int32(v), err
}

func Int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

func Uint(str string) (uint, error) {
	v, err := strconv.ParseUint(str, 10, strconv.IntSize)
	return uint(v), err
}

func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint8(str string) (uint8, error) {
	v, err := strconv.ParseUint(str, 10, 8)
	return uint8(v), err
}

func Uint8ToString(v uint8) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint16(str string) (uint16, error) {
	v, err := strconv.ParseUint(str, 10, 16)
	return uint16(v), err
}

func Uint16ToString(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint32(str string) (uint32, error) {
	v, err := strconv.ParseUint(str, 10, 32)
	return uint32(v), err
}

func Uint32ToString(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Uint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func Float32(str string) (float32, error) {
	v, err := strconv.ParseFloat(str, 32)
	return float32(v), err
}

func Float32ToString(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func Float64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func Bool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

func BoolToString(v bool) string {
	return strconv.FormatBool(v)
}

func String(str string) (string, error) {
	return str, nil
}

func StringToString(v string) string {
	return v
}

func Complex64(str string) (complex64, error) {
	v, err := strconv.ParseComplex(str, 64)
	return complex64(v), err
}

func Complex64ToString(v complex64) string {
	return strconv.FormatComplex(complex128(v), 'f', -1, 64)
}

func Complex128(str string) (complex128, error) {
	return strconv.ParseComplex(str, 128)
}

func Complex128ToString(v complex128) string {
	return strconv.FormatComplex(v, 'f', -1, 128)
}

func Duration(str string) (time.Duration, error) {
	return time.ParseDuration(str)
}

func DurationToString(v time.Duration) string {
	return v.String()
}

func Time(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

func TimeToString(v time.Time) string {
	return v.Format(time.RFC3339)
}

func Location(str string) (time.Location, error) {
	loc, err := time.LoadLocation(str)
	if err != nil {
		loc = time.UTC
	}
	return *loc, err
}

func LocationToString(v time.Location) string {
	return v.String()
}

func Language(str string) (language.Tag, error) {
	return language.Parse(str)
}

func LanguageToString(v language.Tag) string {
	return v.String()
}

func HexByte(str string) (types.HexByte, error) {
	v, err := strconv.ParseUint(str, 16, 8)
	return types.HexByte(v), err
}

func HexByteToString(v types.HexByte) string {
	return strconv.FormatUint(uint64(v), 16)
}

func OctByte(str string) (types.OctByte, error) {
	v, err := strconv.ParseUint(str, 8, 8)
	return types.OctByte(v), err
}

func OctByteToString(v types.OctByte) string {
	return strconv.FormatUint(uint64(v), 8)
}

func HexBytes(str string) (types.HexBytes, error) {
	data, err := hex.DecodeString(str)
	return types.HexBytes(data), err
}

func HexBytesToString(v types.HexBytes) string {
	return hex.EncodeToString([]byte(v))
}

func B32Bytes(str string) (types.B32Bytes, error) {
	data, err := base32.StdEncoding.DecodeString(str)
	return types.B32Bytes(data), err
}

func B32BytesToString(v types.B32Bytes) string {
	return base32.StdEncoding.EncodeToString([]byte(v))
}

func B64Bytes(str string) (types.B64Bytes, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	return types.B64Bytes(data), err
}

func B64BytesToString(v types.B64Bytes) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func Must[P Parseable](parser func(string) (P, error)) func(string) P {
	return func(str string) P {
		v, err := parser(str)
		if err != nil {
			panic(err)
		}
		return v
	}
}

func GetParser[P Parseable](take *P) func(string) (P, error) {
	var p any
	switch any(take).(type) {
	case *int:
		p = any(Int)
	case *int8:
		p = any(Int8)
	case *int16:
		p = any(Int16)
	case *int32:
		p = any(Int32)
	case *int64:
		p = any(Int64)
	case *uint:
		p = any(Uint)
	case *uint8:
		p = any(Uint8)
	case *uint16:
		p = any(Uint16)
	case *uint32:
		p = any(Uint32)
	case *uint64:
		p = any(Uint64)
	case *float32:
		p = any(Float32)
	case *float64:
		p = any(Float64)
	case *bool:
		p = any(Bool)
	case *string:
		p = any(String)
	case *complex64:
		p = any(Complex64)
	case *complex128:
		p = any(Complex128)
	case *time.Duration:
		p = any(Duration)
	case *time.Time:
		p = any(Time)
	case *time.Location:
		p = any(Location)
	case *language.Tag:
		p = any(Language)
	case *[]int:
		p = any(ParseArray(Int))
	case *[]int8:
		p = any(ParseArray(Int8))
	case *[]int16:
		p = any(ParseArray(Int16))
	case *[]int32:
		p = any(ParseArray(Int32))
	case *[]int64:
		p = any(ParseArray(Int64))
	case *[]uint:
		p = any(ParseArray(Uint))
	case *[]uint8:
		p = any(ParseArray(Uint8))
	case *[]uint16:
		p = any(ParseArray(Uint16))
	case *[]uint32:
		p = any(ParseArray(Uint32))
	case *[]uint64:
		p = any(ParseArray(Uint64))
	case *[]float32:
		p = any(ParseArray(Float32))
	case *[]float64:
		p = any(ParseArray(Float64))
	case *[]bool:
		p = any(ParseArray(Bool))
	case *[]string:
		p = any(ParseArray(String))
	case *[]complex64:
		p = any(ParseArray(Complex64))
	case *[]complex128:
		p = any(ParseArray(Complex128))
	case *[]time.Duration:
		p = any(ParseArray(Duration))
	case *[]time.Time:
		p = any(ParseArray(Time))
	case *[]time.Location:
		p = any(ParseArray(Location))
	case *[]language.Tag:
		p = any(ParseArray(Language))
	case *types.HexByte:
		p = any(HexByte)
	case *types.OctByte:
		p = any(OctByte)
	case *types.HexBytes:
		p = any(HexBytes)
	case *types.B32Bytes:
		p = any(B32Bytes)
	case *types.B64Bytes:
		p = any(B64Bytes)
	case *[]types.HexByte:
		p = any(ParseArray(HexByte))
	case *[]types.OctByte:
		p = any(ParseArray(OctByte))
	case *[]types.HexBytes:
		p = any(ParseArray(HexBytes))
	case *[]types.B32Bytes:
		p = any(ParseArray(B32Bytes))
	case *[]types.B64Bytes:
		p = any(ParseArray(B64Bytes))
	}
	return p.(func(string) (P, error))
}

func Parse[P Parseable](take *P, str string) error {
	parse := GetParser(take)
	var err error
	*take, err = parse(str)
	return err
}

func MustParse[P Parseable](take *P, str string) {
	if err := Parse(take, str); err != nil {
		panic(err)
	}
}

func ParseArray[P Parseable](parser func(string) (P, error)) func(string) ([]P, error) {
	return func(str string) ([]P, error) {
		ret := []P{}
		var part string
		var v P
		var err error
		parts := strings.Split(str, ",")
		for i := range parts {
			if part = strings.TrimSpace(parts[i]); part != "" {
				if v, err = parser(part); err != nil {
					return nil, err
				}
				ret = append(ret, v)
			}
		}
		return ret, nil
	}
}

func ArrayToString[P Parseable](encoder func(P) string) func([]P) string {
	return func(v []P) string {
		if len(v) == 0 {
			return ""
		}
		ret := ""
		for i := 0; i < len(v)-1; i++ {
			ret = ret + encoder(v[i]) + ","
		}
		ret = ret + encoder(v[len(v)-1])
		return ret
	}
}

func ToString[P Parseable](v P) string {
	switch p := any(v).(type) {
	case int:
		return IntToString(p)
	case int8:
		return Int8ToString(p)
	case int16:
		return Int16ToString(p)
	case int32:
		return Int32ToString(p)
	case int64:
		return Int64ToString(p)
	case uint:
		return UintToString(p)
	case uint8:
		return Uint8ToString(p)
	case uint16:
		return Uint16ToString(p)
	case uint32:
		return Uint32ToString(p)
	case uint64:
		return Uint64ToString(p)
	case float32:
		return Float32ToString(p)
	case float64:
		return Float64ToString(p)
	case bool:
		return BoolToString(p)
	case string:
		return StringToString(p)
	case complex64:
		return Complex64ToString(p)
	case complex128:
		return Complex128ToString(p)
	case time.Duration:
		return DurationToString(p)
	case time.Time:
		return TimeToString(p)
	case time.Location:
		return LocationToString(p)
	case language.Tag:
		return LanguageToString(p)
	case []int:
		return ArrayToString(IntToString)(p)
	case []int8:
		return ArrayToString(Int8ToString)(p)
	case []int16:
		return ArrayToString(Int16ToString)(p)
	case []int32:
		return ArrayToString(Int32ToString)(p)
	case []int64:
		return ArrayToString(Int64ToString)(p)
	case []uint:
		return ArrayToString(UintToString)(p)
	case []uint8:
		return ArrayToString(Uint8ToString)(p)
	case []uint16:
		return ArrayToString(Uint16ToString)(p)
	case []uint32:
		return ArrayToString(Uint32ToString)(p)
	case []uint64:
		return ArrayToString(Uint64ToString)(p)
	case []float32:
		return ArrayToString(Float32ToString)(p)
	case []float64:
		return ArrayToString(Float64ToString)(p)
	case []bool:
		return ArrayToString(BoolToString)(p)
	case []string:
		return ArrayToString(StringToString)(p)
	case []complex64:
		return ArrayToString(Complex64ToString)(p)
	case []complex128:
		return ArrayToString(Complex128ToString)(p)
	case []time.Duration:
		return ArrayToString(DurationToString)(p)
	case []time.Time:
		return ArrayToString(TimeToString)(p)
	case []time.Location:
		return ArrayToString(LocationToString)(p)
	case []language.Tag:
		return ArrayToString(LanguageToString)(p)
	case types.HexByte:
		return HexByteToString(p)
	case types.OctByte:
		return OctByteToString(p)
	case types.HexBytes:
		return HexBytesToString(p)
	case types.B32Bytes:
		return B32BytesToString(p)
	case types.B64Bytes:
		return B64BytesToString(p)
	case []types.HexByte:
		return ArrayToString(HexByteToString)(p)
	case []types.OctByte:
		return ArrayToString(OctByteToString)(p)
	case []types.HexBytes:
		return ArrayToString(HexBytesToString)(p)
	case []types.B32Bytes:
		return ArrayToString(B32BytesToString)(p)
	case []types.B64Bytes:
		return ArrayToString(B64BytesToString)(p)
	}
	panic("should not reach")
}
