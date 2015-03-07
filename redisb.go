package redisb

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const nilLength = -1
type RedisNil struct{}

type InternalError struct {
	e error
}

func (e InternalError) Error() string {
	return e.Error()
}

type RedisError struct {
	prefix string
	suffix string
}

func (e RedisError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.prefix, e.suffix)
}

func parseError(s string) RedisError {
	p := strings.SplitN(s, " ", 2)
	if len(p) < 2 {
		p[1] = "N/A"
	}
	return RedisError{p[0], p[1]}
}

func Do(c net.Conn, args ...string) (interface{}, error) {
	fmt.Fprintf(c, encode(args))
	return decode(bufio.NewReader(c))
}

func encode(i interface{}) string {
	switch t := i.(type) {
	case []string:
		s := []string{"*", strconv.Itoa(len(i.([]string))), "\r\n"}
		for _, v := range i.([]string) {
			s = append(s, encode(v))
		}
		return strings.Join(s, "")
	case string:
		return "$" + strconv.Itoa(len(i.(string))) + "\r\n" + i.(string) + "\r\n"
	default:
		panic(fmt.Sprintf("Unable to encode type: %#v", t))
	}
}

func decode(r *bufio.Reader) (interface{}, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, InternalError{err}
	}
	switch string(t) {
	case "-":
		s, err := r.ReadString('\n')
		if err != nil {
			return nil, InternalError{err}
		}
		return nil, parseError(s)
	case "+":
		return r.ReadString('\n')
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
	s, err := r.ReadString('\n')
	if err != nil {
		return nil, InternalError{err}
	}
	i, err := toInt(s)
	if err != nil {
		return nil, InternalError{err}
	}
	return i, nil
}

func decodeBulkStringSuffix(r *bufio.Reader) (interface{}, error) {
	tmp, err := r.ReadString('\n')
	if err != nil {
		return nil, InternalError{err}
	}
	slen, err := toInt(tmp)
	if err != nil {
		return nil, InternalError{err}
	}
	if slen == nilLength {
		return RedisNil{}, nil
	}
	s := make([]byte, slen)
	r.Read(s)
	r.ReadByte()
	r.ReadByte()
	return string(s), nil
}

func decodeArraySuffix(r *bufio.Reader) (interface{}, error) {
	tmp, err := r.ReadString('\n')
	if err != nil {
		return nil, InternalError{err}
	}
	alen, err := toInt(tmp)
	if err != nil {
		return nil, InternalError{err}
	}
	if alen == nilLength {
		return RedisNil{}, nil
	}
	result := make([]interface{}, 0, alen)
	for i := int64(0); i < alen; i++ {
		v, err := decode(r)
		if err != nil {
			return nil, InternalError{err}
		}
		result = append(result, v)
	}
	return result, nil
}

func toInt(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}
