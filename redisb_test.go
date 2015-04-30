package redisb

import (
	"testing"
)

/*
type RedisNil struct{}
func Do(c net.Conn, args ...string) (interface{}, error) {
func DoN(c net.Conn, args ...string) (interface{}, error) {
func Out(c net.Conn, args ...string) {
func encode(i interface{}) string {
func decode(r *bufio.Reader) (interface{}, error) {
func decodeIntSuffix(r *bufio.Reader) (interface{}, error) {
func decodeBulkStringSuffix(r *bufio.Reader) (interface{}, error) {
func decodeArraySuffix(r *bufio.Reader) (interface{}, error) {
func toInt(s string) (int64, error) {
*/

func TestParseError(t *testing.T) {
        cases := []struct {
                in string
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

func TestToInt(t *testing.T) {
        cases := []struct {
                in string
                out struct {
                        r int64
                        e error
                }
        }{
                {"1", struct{int64, error}{1, nil}},
        }
        for _, c := range cases {
                got := toInt(c.in)
                if got != c.out {
                        t.Errorf("toInt: %s", c.in)
                }
        }
}
