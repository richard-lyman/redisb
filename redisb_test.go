package redisb

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	c, err := net.Dial("tcp", "localhost:6379")
	if err == nil {
		Do(c, "GET", "some_key_that_should_not_have_a_value")
		Do(c, "SET", "a", "b")
	}
}

func TestDoN(t *testing.T) {
	c, err := net.Dial("tcp", "localhost:6379")
	if err == nil {
		DoN(c, "GET", "some_key_that_should_not_have_a_value")
		DoN(c, "SET", "a", "b")
	}
}

func TestOut(t *testing.T) {
	c, err := net.Dial("tcp", "localhost:6379")
	if err == nil {
		Out(c, "GET", "some_key_that_should_not_have_a_value")
		Out(c, "SET", "a", "b")
	}
}

func TestErrorConstructors(t *testing.T) {
	if newConnError("a").Error() != "a" {
		t.Error("newConnError failed to report Error correctly")
	}
	if newConversionError("a").Error() != "a" {
		t.Error("newConversionError failed to report Error correctly")
	}
	if parseError("a b").Error() != "[a]: b" {
		t.Error("parseError failed to report Error correctly")
	}
}

func TestDecode(t *testing.T) {
	bs := func(s string) *bufio.Reader { return bufio.NewReader(strings.NewReader(s)) }
	cases := []struct {
		in  string
		out interface{}
		err error
	}{
		{"", nil, newConnError("")},
		// Simple strings
		{"+", nil, newConnError("")},
		{"+a\r\n", "a", nil},
		// Errors
		{"-ERROR_TYPE Some error message\r\n", nil, parseError("ERROR_TYPE Some error message")},
		// Ints
		{":", nil, newConnError("")},
		{":1\r\n", 1, nil},
		{":a\r\n", nil, newConversionError("")},
		// Bulk strings
		{"$", nil, newConnError("")},
		{"$a\r\n", nil, newConversionError("")},
		{"$1\r\na\r\n", "a", nil},
		{"$1\r\n", nil, errors.New("")},
		{"$2\r\na\r\n", nil, errors.New("")},
		{"$-1\r\n", RedisNil{}, nil},
		// Arrays
		{"*", nil, newConnError("")},
		{"*0\r\n", []string{}, nil},
		{"*-1\r\n", RedisNil{}, nil},
		{"*1\r\n:1\r\n", []int{1}, nil},
		{"*1\r\n", nil, errors.New("")},
		{"*2\r\n:1\r\n", nil, errors.New("")},
	}
	for _, c := range cases {
		tmp, err := decode(bs(c.in))
		if c.err == nil && err != nil {
			t.Errorf("decode error expectations not met: %q, %s, %s", c.in, c.err.Error(), err.Error())
		}
		fmt.Printf("decode: %q: %v - %v\n", c.in, c.out, tmp)
	}
}

func TestToInt(t *testing.T) {
	cases := []struct {
		in  string
		out int64
		err error
	}{
		{"1", int64(1), nil},
		{"10", int64(10), nil},
		{"-1", int64(-1), nil},
		{"a", int64(0), errors.New(`strconv.ParseInt: parsing "a": invalid syntax`)},
		{"", int64(0), errors.New(`strconv.ParseInt: parsing "": invalid syntax`)},
	}
	for _, c := range cases {
		tmp, err := toInt(c.in)
		if err != nil && err.Error() != c.err.Error() {
			t.Errorf("toInt error expectations not met: %q, %s, %s", c.in, c.err.Error(), err.Error())
		}
		if tmp != c.out {
			t.Errorf("toInt: %s: %q - %q", c.in, c.out, tmp)
		}
	}
}

func TestCleanEnding(t *testing.T) {
	if cleanEnding("a") != "a" {
		t.Error("cleanEnding failed to keep clean input clean")
	}
	if cleanEnding("a\r\n") != "a" {
		t.Error("cleanEnding failed to clean input that needed cleaning")
	}
}

func TestEncodePanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Encode failed to panic")
		}
	}()
	encode(1)
}

func TestEncode(t *testing.T) {
	cases := []struct {
		in  interface{}
		out string
	}{
		{"a", "$1\r\na\r\n"},
		{[]string{"a"}, "*1\r\n$1\r\na\r\n"},
		{[]string{"a", "b"}, "*2\r\n$1\r\na\r\n$1\r\nb\r\n"},
	}
	for _, c := range cases {
		tmp := encode(c.in)
		if tmp != c.out {
			t.Errorf("encode: %s: %q - %q", c.in, c.out, tmp)
		}
	}
}

func TestEncodePanis(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Encode failed to panic when given unacceptable input")
		}
	}()
	unencodableInput := 1
	encode(unencodableInput)
}

func TestParseError(t *testing.T) {
	cases := []struct {
		in  string
		out RedisError
	}{
		{"k v", RedisError{"k", "v"}},
		{"k", RedisError{"k", "N/A"}},
		{"", RedisError{"", "N/A"}},
	}
	for _, c := range cases {
		got := parseError(c.in)
		if got != c.out {
			t.Errorf("parseInt: %s", c.in)
		}
	}
}
