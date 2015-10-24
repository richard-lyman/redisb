package redisb

import (
	"errors"
	"testing"
)

/*
type RedisNil struct{}
func Do(c net.Conn, args ...string) (interface{}, error) {
func DoN(c net.Conn, args ...string) (interface{}, error) {
func Out(c net.Conn, args ...string) {
func decode(r *bufio.Reader) (interface{}, error) {
func decodeIntSuffix(r *bufio.Reader) (interface{}, error) {
func decodeBulkStringSuffix(r *bufio.Reader) (interface{}, error) {
func decodeArraySuffix(r *bufio.Reader) (interface{}, error) {
func toInt(s string) (int64, error) {
*/

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
