package redisb

import (
	"bufio"
	"strconv"
	//"bytes"
	"errors"
	//"fmt"
	"io/ioutil"
	"net"
	"strings"
	"testing"
)

func benchmarkEncodeDirectN(b *testing.B, n int) {
	w := bufio.NewWriterSize(ioutil.Discard, 10000)
	v := []string{strings.Repeat("b", n)}
	/*
		v := []string{}
		for i := 0; i < n; i++ {
			v = append(v, "b")
		}
	*/
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodeDirect(v, w)
		w.Flush()
	}
}

func BenchmarkEncodeDirect10(b *testing.B)    { benchmarkEncodeDirectN(b, 10) }
func BenchmarkEncodeDirect100(b *testing.B)   { benchmarkEncodeDirectN(b, 100) }
func BenchmarkEncodeDirect1000(b *testing.B)  { benchmarkEncodeDirectN(b, 1000) }
func BenchmarkEncodeDirect10000(b *testing.B) { benchmarkEncodeDirectN(b, 10000) }

func benchmarkEncodeN(b *testing.B, n int, m int, f func(interface{}) string) {
	var r string
	//v := strings.Repeat("b", n)
	v := []string{}
	for i := 0; i < n; i++ {
		v = append(v, strings.Repeat("b", m))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = f(v)
	}
	result = r
}

func BenchmarkEncode111(b *testing.B)     { benchmarkEncodeN(b, 1, 1, encode) }
func BenchmarkEncode112(b *testing.B)     { benchmarkEncodeN(b, 1, 10, encode) }
func BenchmarkEncode113(b *testing.B)     { benchmarkEncodeN(b, 1, 100, encode) }
func BenchmarkEncode114(b *testing.B)     { benchmarkEncodeN(b, 1, 1000, encode) }
func BenchmarkEncode1101(b *testing.B)    { benchmarkEncodeN(b, 10, 10, encode) }
func BenchmarkEncode1102(b *testing.B)    { benchmarkEncodeN(b, 10, 100, encode) }
func BenchmarkEncode1103(b *testing.B)    { benchmarkEncodeN(b, 10, 1000, encode) }
func BenchmarkEncode11001(b *testing.B)   { benchmarkEncodeN(b, 100, 10, encode) }
func BenchmarkEncode11002(b *testing.B)   { benchmarkEncodeN(b, 100, 100, encode) }
func BenchmarkEncode11003(b *testing.B)   { benchmarkEncodeN(b, 100, 1000, encode) }
func BenchmarkEncode110001(b *testing.B)  { benchmarkEncodeN(b, 1000, 10, encode) }
func BenchmarkEncode110002(b *testing.B)  { benchmarkEncodeN(b, 1000, 100, encode) }
func BenchmarkEncode110003(b *testing.B)  { benchmarkEncodeN(b, 1000, 1000, encode) }
func BenchmarkEncode1100001(b *testing.B) { benchmarkEncodeN(b, 10000, 10, encode) }
func BenchmarkEncode1100002(b *testing.B) { benchmarkEncodeN(b, 10000, 100, encode) }
func BenchmarkEncode1100003(b *testing.B) { benchmarkEncodeN(b, 10000, 1000, encode) }

func BenchmarkEncode211(b *testing.B)     { benchmarkEncodeN(b, 1, 1, encode2) }
func BenchmarkEncode212(b *testing.B)     { benchmarkEncodeN(b, 1, 10, encode2) }
func BenchmarkEncode213(b *testing.B)     { benchmarkEncodeN(b, 1, 100, encode2) }
func BenchmarkEncode214(b *testing.B)     { benchmarkEncodeN(b, 1, 1000, encode2) }
func BenchmarkEncode2101(b *testing.B)    { benchmarkEncodeN(b, 10, 10, encode2) }
func BenchmarkEncode2102(b *testing.B)    { benchmarkEncodeN(b, 10, 100, encode2) }
func BenchmarkEncode2103(b *testing.B)    { benchmarkEncodeN(b, 10, 1000, encode2) }
func BenchmarkEncode21001(b *testing.B)   { benchmarkEncodeN(b, 100, 10, encode2) }
func BenchmarkEncode21002(b *testing.B)   { benchmarkEncodeN(b, 100, 100, encode2) }
func BenchmarkEncode21003(b *testing.B)   { benchmarkEncodeN(b, 100, 1000, encode2) }
func BenchmarkEncode210001(b *testing.B)  { benchmarkEncodeN(b, 1000, 10, encode2) }
func BenchmarkEncode210002(b *testing.B)  { benchmarkEncodeN(b, 1000, 100, encode2) }
func BenchmarkEncode210003(b *testing.B)  { benchmarkEncodeN(b, 1000, 1000, encode2) }
func BenchmarkEncode2100001(b *testing.B) { benchmarkEncodeN(b, 10000, 10, encode2) }
func BenchmarkEncode2100002(b *testing.B) { benchmarkEncodeN(b, 10000, 100, encode2) }
func BenchmarkEncode2100003(b *testing.B) { benchmarkEncodeN(b, 10000, 1000, encode2) }

func BenchmarkParseInt(b *testing.B) {
	var r int64
	for i := 0; i < b.N; i++ {
		r, _ = strconv.ParseInt("123", 10, 64)
	}
	result = r

}

func BenchmarkParseUint(b *testing.B) {
	var r uint64
	for i := 0; i < b.N; i++ {
		r, _ = strconv.ParseUint("123", 10, 64)
	}
	result = r

}

var result interface{}

func benchmarkDoGetN(b *testing.B, n int) {
	var r interface{}
	c, err := net.Dial("tcp", "localhost:6379")
	if err == nil {
		//v := strings.Repeat("b", n)
		//Do(c, "SET", "benchdoget", v)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r, _ = Do(c, "GET", "benchdoget")
		}
	} else {
		b.Errorf("Failed to connect: %s", err)
	}
	result = r
}

func BenchmarkDoGet10(b *testing.B)        { benchmarkDoGetN(b, 10) }
func BenchmarkDoGet1000(b *testing.B)      { benchmarkDoGetN(b, 1000) }
func BenchmarkDoGet100000(b *testing.B)    { benchmarkDoGetN(b, 100000) }
func BenchmarkDoGet10000000(b *testing.B)  { benchmarkDoGetN(b, 10000000) }
func BenchmarkDoGet100000000(b *testing.B) { benchmarkDoGetN(b, 100000000) }

func TestDo(t *testing.T) {
	c, err := net.Dial("tcp", "localhost:6379")
	if err == nil {
		Do(c, "GET", "some_key_that_should_not_have_a_value")
		Do(c, "SET", "a", "b")
	}
}

func TestDoSubscribe(t *testing.T) {
	c, err := net.Dial("tcp", "localhost:6379")
	if err == nil {
		Do(c, "SUBSCRIBE", "some_channel") // TODO - should this work??
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
		//tmp, err := decode(bs(c.in))
		_, err := decode(bs(c.in))
		if c.err == nil && err != nil {
			t.Errorf("decode error expectations not met: %q, %s", c.in, err.Error())
		}
		//fmt.Printf("decode: %q: %v - %v\n", c.in, c.out, tmp)
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

func TestEncode(t *testing.T) {
	cases := []struct {
		in  []string
		out string
	}{
		{[]string{"a"}, "*1\r\n$1\r\na\r\n"},
		{[]string{"a", "b"}, "*2\r\n$1\r\na\r\n$1\r\nb\r\n"},
	}
	for _, c := range cases {
		//var b bytes.Buffer
		//encode(c.in, &b)
		//tmp := b.String()
		tmp := encode(c.in)
		if tmp != c.out {
			t.Errorf("encode: %s: %q - %q", c.in, c.out, tmp)
		}
	}
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
