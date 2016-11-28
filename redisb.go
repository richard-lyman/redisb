package redisb

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type redisType string

var returnTypes = map[string]redisType{
	"set": "bool",
}

func Raw(rw io.ReadWriter, args ...string) (interface{}, error) {
	fmt.Fprint(rw, Encode(args))
	return Decode(bufio.NewReader(rw))
}

func Int64(rw io.ReadWriter, args ...string) (int64, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return 0, err
	}
	result, err := toInt64(i)
	return result, err
}

func toInt64(i interface{}) (int64, error) {
	switch t := i.(type) {
	case int64:
		return t, nil
	case string:
		result, err := toInt(t)
		if err != nil {
			return 0, newConversionError("Conversion to int64 failed: %#v, %s", i, err)
		}
		return result, nil
	}
	return 0, newConversionError("Conversion to int64 failed: %#v", i)
}

func Bool(rw io.ReadWriter, args ...string) (bool, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return false, err
	}
	result, err := toBool(i)
	return result, err
}

func toBool(i interface{}) (bool, error) {
	switch i {
	case "OK":
		return true, nil
	case "1":
		return true, nil
	case int64(1):
		return true, nil
	case "0":
		return false, nil
	case int64(0):
		return false, nil
	case nil:
		return false, nil
	}
	return false, newConversionError("Conversion to bool failed: %#v %s", i, reflect.TypeOf(i))
}

func String(rw io.ReadWriter, args ...string) (string, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return "", err
	}
	result, err := toString(i)
	return result, err
}

func toString(i interface{}) (string, error) {
	s, ok := i.(string)
	if ok {
		return s, nil
	}
	return "", newConversionError("Conversion to string failed: %#v", i)
}

func Array(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return nil, err
	}
	a, ok := i.([]interface{})
	if ok {
		return a, nil
	}
	return nil, newConversionError("Conversion to []interface{} failed: %#v", i)
}

func Bools(rw io.ReadWriter, args ...string) ([]bool, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return nil, err
	}
	a, ok := i.([]interface{})
	if !ok {
		return nil, newConversionError("Conversion to []bool failed: %#v", i)
	}
	result := []bool{}
	for _, v := range a {
		sv, err := toBool(v)
		if err != nil {
			return nil, newConversionError("Conversion to []bool failed: %#v: %s", i, err)
		}
		result = append(result, sv)
	}
	return result, nil
}

func Int64s(rw io.ReadWriter, args ...string) ([]int64, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return nil, err
	}
	a, ok := i.([]interface{})
	if !ok {
		return nil, newConversionError("Conversion to []int64 failed: %#v", i)
	}
	result := []int64{}
	for _, v := range a {
		sv, err := toInt64(v)
		if err != nil {
			return nil, newConversionError("Conversion to []int64 failed: %#v: %s", i, err)
		}
		result = append(result, sv)
	}
	return result, nil
}

func Strings(rw io.ReadWriter, args ...string) ([]string, error) {
	fmt.Fprint(rw, Encode(args))
	i, err := Decode(bufio.NewReader(rw))
	if err != nil {
		return nil, err
	}
	a, ok := i.([]interface{})
	if !ok {
		return nil, newConversionError("Conversion to []string failed: %#v", i)
	}
	result := []string{}
	for _, v := range a {
		sv, err := toString(v)
		if err != nil {
			return nil, newConversionError("Conversion to []string failed: %#v: %s", i, err)
		}
		result = append(result, sv)
	}
	return result, nil
}

type ReaderError struct {
	e error
}

func (re ReaderError) Error() string {
	return re.e.Error()
}

func newReaderError(format string, values ...interface{}) ReaderError {
	return ReaderError{fmt.Errorf(format, values...)}
}

type ConversionError struct {
	e error
}

func (ce ConversionError) Error() string {
	return ce.e.Error()
}

func newConversionError(format string, values ...interface{}) ConversionError {
	return ConversionError{fmt.Errorf(format, values...)}
}

type RedisError struct {
	Prefix string
	Suffix string
}

func (re RedisError) Error() string {
	return fmt.Sprintf("[%s]: %s", re.Prefix, re.Suffix)
}

func parseError(s string) RedisError {
	p := strings.SplitN(s, " ", 2)
	if len(p) < 2 {
		p = append(p, "N/A")
	}
	return RedisError{p[0], p[1]}
}

func Encode(i interface{}) string {
	switch t := i.(type) {
	case []string:
		s := []string{"*", strconv.Itoa(len(t)), "\r\n"}
		for _, v := range t {
			s = append(s, Encode(v))
		}
		return strings.Join(s, "")
	case string:
		return "$" + strconv.Itoa(len(t)) + "\r\n" + t + "\r\n"
	default:
		panic(fmt.Sprintf("Unable to Encode type: %#v", t))
	}
}

func Decode(r *bufio.Reader) (interface{}, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, newReaderError("Failed to get Redis type byte in to call ReadByte: %s", err)
	}
	//fmt.Println("Type:", string(t))
	switch string(t) {
	case "-":
		s, err := redisReadString(r)
		if err != nil {
			return nil, newReaderError("Failed to get Error string in call to ReadString: %s", err)
		}
		return nil, parseError(s)
	case "+":
		tmp, err := redisReadString(r)
		return tmp, err
	case ":":
		return decodeIntSuffix(r)
	case "$":
		return decodeBulkStringSuffix(r)
	case "*":
		return decodeArraySuffix(r)
	}
	panic(fmt.Sprintf("Failed to identify type: '%q'", string(t)))
}

func decodeIntSuffix(r *bufio.Reader) (interface{}, error) {
	s, err := redisReadString(r)
	if err != nil {
		return nil, newReaderError("Failed to get raw int in call to ReadString: %s", err)
	}
	i, err := toInt(s)
	if err != nil {
		return nil, newConversionError("Failed to convert raw int to int: %s", err)
	}
	return i, nil
}

func decodeBulkStringSuffix(r *bufio.Reader) (interface{}, error) {
	tmp, err := redisReadString(r)
	if err != nil {
		return nil, newReaderError("Failed to get raw int for Bulk String size in call to ReadString: %s", err)
	}
	if isNegativeOne(tmp) {
		//fmt.Println("Negative one - redis null on bulk empty string")
		return nil, nil
	}
	slen, err := toUint(tmp)
	if err != nil {
		return nil, newConversionError("Failed to convert raw int to int for Bulk String size: %s", err)
	}
	s := make([]byte, slen)
	_, err = io.ReadFull(r, s)
	if err == io.EOF {
		return nil, fmt.Errorf("Unable to read any bytes")
	}
	if err != nil {
		return nil, fmt.Errorf("Unable to read required number of bytes: %s", err)
	}
	r.ReadByte()
	r.ReadByte()
	return string(s), nil
}

func decodeArraySuffix(r *bufio.Reader) (interface{}, error) {
	tmp, err := redisReadString(r)
	if err != nil {
		return nil, newReaderError("Failed to get raw int for Bulk Array size in call to ReadString: %s", err)
	}
	if isNegativeOne(tmp) {
		return nil, nil
	}
	alen, err := toUint(tmp)
	if err != nil {
		return nil, newConversionError("Failed to convert raw int to int for Bulk Array size: %s", err)
	}
	result := make([]interface{}, 0, alen)
	for i := uint64(0); i < alen; i++ {
		v, err := Decode(r)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func isNegativeOne(s string) bool {
	return len(s) == 2 && s[0] == '-' && s[1] == '1'
}

func toUint(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(s), 10, 64)
}

func toInt(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}

func redisReadString(r *bufio.Reader) (string, error) {
	var out bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", fmt.Errorf("failed to read byte: %s", err)
		}
		if b == '\r' {
			b, err := r.ReadByte()
			if err != nil {
				return "", fmt.Errorf("failed to read byte: %s", err)
			}
			if b != '\n' {
				return "", fmt.Errorf("failed to read required final newline byte")
			}
			//return strings.TrimSuffix(out.String(), "\r\n"), nil
			return out.String(), nil
		}
		out.WriteByte(b)
	}
}
