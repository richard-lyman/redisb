package redisb

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type RedisNil struct{}

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

func decode(r *bufio.Reader) (interface{}, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch string(t) {
	case "-":
		s, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return nil, parseError(s)
	case "+":
		return r.ReadString('\n')
	case ":":
		s, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		i, err := toInt(s)
		if err != nil {
			return nil, err
		}
		return i, nil
	case "$":
		tmp, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		slen, err := toInt(tmp)
		if err != nil {
			return nil, err
		}
		if slen == -1 {
			return RedisNil{}, nil
		}
		s := make([]byte, slen)
		r.Read(s)
		r.ReadByte()
		r.ReadByte()
		return string(s), nil
	case "*":
		tmp, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		alen, err := toInt(tmp)
		if err != nil {
			return nil, err
		}
		if alen == -1 {
			return RedisNil{}, nil
		}
		result := make([]interface{}, 0, alen)
		for i := int64(0); i < alen; i++ {
			v, err := decode(r)
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}
		return result, nil
	}
	panic(fmt.Sprintf("Failed to identify type: '%q'", string(t)))
}

func toInt(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
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
