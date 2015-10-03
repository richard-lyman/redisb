/*
Package redisb implements a simple Redis Client - a 'Redis base' or redisb.

For any Redis command one option is to use Do to send the request and process the response as follows:

        c, err := net.Dial("tcp", "localhost:6379")
        if err != nil {
                panic(err.Error())
        }
        tmpv, err := redisb.Do(c, "GET", "some_key")
        if err != nil {
                // Handle the error, possibly as shown below with a type switch
        } else if tmpv == nil {
                // Handle a nil value - that occurred without error
        } else {
                // There's either an int64, a string, or an interface{} slice in tmpv depending on the Redis Command used
        }

For handling errors, you might use a type switch as follows:

        switch err.(type) {
        case redisb.RedisError:
                // Redis saw your request, and there is a related error provided from Redis
        case redisb.ConnError:
                // The given net.Conn is 'bad' - if it came from a pool, just Close it - don't return it.
                // You might be successful if you opened a new net.Conn and tried that same call again.
                // The Redis Server you connected to might not be valid (unreachable, etc.). You might be successful if you opened a net.Conn to another Server.
        case redisb.ConversionError:
                // There is something seriously wrong...
                // ... your connection is good, Redis saw your request, Redis didn't have a problem responding... and it gave an invalid response
                // Submit a bug. :-)
        default:
                // This is some other error - which shouldn't occur... again, submit a bug. :-)
        }
*/
package redisb

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const nilLength = -1

// RedisNil represents the unique Redis Null value that can result from Bulk String or Bulk Array.
// Signals, 'there isn't a value,' which is different than an empty Bulk String value or empty Bulk Array value.
type RedisNil struct{}

// ConnError indicates a failure to read from the given net.Conn.
type ConnError struct {
	e error
}

func (ce ConnError) Error() string {
	return ce.e.Error()
}

func newConnError(format string, values ...interface{}) ConnError {
	return ConnError{fmt.Errorf(format, values...)}
}

// ConversionError indicates a failure to convert from the raw bytes representing an int to an int.
type ConversionError struct {
	e error
}

func (ce ConversionError) Error() string {
	return ce.e.Error()
}

func newConversionError(format string, values ...interface{}) ConversionError {
	return ConversionError{fmt.Errorf(format, values...)}
}

// RedisError wraps the Redis Error data type.
type RedisError struct {
	// Holds the prefix key provided by Redis to indicate what class of errors the Suffix describes
	Prefix string
	// Holds the suffix description provided by Redis to describe this specific instance of the class of error indicated by the Prefix
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

/*
Do accepts a net.Conn, connected to an instance of Redis, and any number of strings that are the Bulk String values in the Bulk Array formatted request sent to Redis.
Do will return the response from Redis.
The actual type of the value returned depends on the Redis command used in the request.

Do internally calls DoN and converts the RedisNil value returned from DoN to a Go nil so it's easier to handle.
The errors returned have the same possibilities as DoN.

Outside of the normal responses, there is one other unique response that can be returned.
Redis can reply with a value that holds an error message.
That error message is encoded in the error return value as a redisb.RedisError.

Do can be used as a Send/Receive or as just a Receive.
When len(args) > 0, they are written to the net.Conn and then the response is received.
When len(args) == 0, then nothing is written to the net.Conn and the response is received.
This allows Do to work for notifications.

For any Redis command one option to process a Do response is as follows:

        c, err := net.Dial("tcp", "localhost:6379")
        if err != nil {
                panic(err.Error())
        }
        tmpv, err := redisb.Do(c, "GET", "some_key")
        if err != nil {
                // Handle the error, possibly as shown below with a type switch
        } else if tmpv == nil {
                // Handle a nil value - that occurred without error
        } else {
                // There's either an int64, a string, or an interface{} slice in tmpv depending on the Redis Command used
        }

For handling errors, you might use a type switch as follows:

        switch err.(type) {
        case redisb.RedisError:
                // Redis saw your request, and there is a related error provided from Redis
        case redisb.ConnError:
                // The given net.Conn is 'bad' - if it came from a pool, just Close it - don't return it.
                // You might be successful if you opened a new net.Conn and tried that same call again.
                // The Redis Server you connected to might not be valid (unreachable, etc.). You might be successful if you opened a net.Conn to another Server.
        case redisb.ConversionError:
                // There is something seriously wrong...
                // ... your connection is good, Redis saw your request, Redis didn't have a problem responding... and it gave an invalid response
                // Submit a bug. :-)
        default:
                // This is some other error - which shouldn't occur... again, submit a bug. :-)
        }
*/
func Do(c net.Conn, args ...string) (interface{}, error) {
	tmp, err := DoN(c, args...)
	if _, isRedisNil := tmp.(RedisNil); isRedisNil {
		return nil, err
	}
	return tmp, err
}

/*
DoN accepts a net.Conn, connected to an instance of Redis, and any number of strings that are the Bulk String values in the Bulk Array formatted request sent to Redis.
DoN will return the response from Redis.
The actual type of the value returned depends on the Redis command used in the request.

Outside of the normal responses, there are two other unique values that can be returned.
Redis can reply with a value that holds an error message.
That error message is encoded in the error return value as a redisb.RedisError.
The other unique response value is not returned in the error, and has the redisb.RedisNil type.

DoN can be used as a Send/Receive or as just a Receive.
When len(args) > 0, they are written to the net.Conn and then the response is received.
When len(args) == 0, then nothing is written to the net.Conn and the response is received.
This allows DoN to work for notifications.

For any Redis command one option to process the response is as follows:

        c, err := net.Dial("tcp", "localhost:6379")
        if err != nil {
                panic(err.Error())
        }
        tmpv, err := redisb.DoN(c, "GET", "some_key")
        if _, isRedisNil := tmpv.(redisb.RedisNil); err == nil && !isRedisNil {
                // There's a concrete non-nil value in tmpv
        } else {
                // Either there's an error in err, or we need to handle a RedisNil value

        }

For any Redis command another option to process the response is as follows:

        c, err := net.Dial("tcp", "localhost:6379")
        if err != nil {
                panic(err.Error())
        }
        tmpv, err := redisb.DoN(c, "GET", "some_key")
        if tmpv == nil {
                // Handle the many different error types (possibly as shown below)
        } else {
                // If the 'err' is nil, then there is a sane value in tmpv - just make sure you handle the Redis Nil value
                switch tmpv.(type) {
                case redisb.RedisNil:
                        // Handle
                default:
                        // Convert to a type depending on the Redis command used
                }
        }

For any Redis command another option to process the response is as follows:

        c, err := net.Dial("tcp", "localhost:6379")
        if err != nil {
                panic(err.Error())
        }
        tmpv, err := redisb.DoN(c, "GET", "some_key")
        if err != nil {
                switch err.(type) {
                case redisb.RedisError:
                        // Redis saw your request, and there is a related error provided from Redis
                case redisb.ConnError:
                        // The given net.Conn is 'bad' - if it came from a pool, just Close it - don't return it.
                        // You might be successful if you opened a new net.Conn and tried that same call again.
                        // The Redis Server you connected to might not be valid (unreachable, etc.). You might be successful if you opened a net.Conn to another Server.
                case redisb.ConversionError:
                        // There is something seriously wrong...
                        // ... your connection is good, Redis saw your request, Redis didn't have a problem responding... and it gave an invalid response
                        // Submit a bug. :-)
                default:
                        // This is some other error - which shouldn't occur... again, submit a bug. :-)
                }
        }
        // Use tmpv as a specific value, who's type depends on the Redis Command used.
        // The type of tmpv could be an int64, string, interface{} slice, or the special redisb.RedisNil.
        // You may want to use code like the following:
        switch tmpv.(type) {
        case redisb.RedisNil:
                // Handle
        default:
                // Convert to a type depending on the Redis command used
        }
*/
func DoN(c net.Conn, args ...string) (interface{}, error) {
	if len(args) > 0 {
		fmt.Fprint(c, encode(args))
	}
	return decode(bufio.NewReader(c))
}

// Out sends the arguments and does not read any of the response(s) back in
func Out(c net.Conn, args ...string) {
	fmt.Fprint(c, encode(args))
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

func cleanEnding(s string) string {
	return strings.TrimSuffix(s, "\r\n")
}

func decode(r *bufio.Reader) (interface{}, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, newConnError("Failed to get Redis type byte in to call ReadByte: %s", err)
	}
	switch string(t) {
	case "-":
		s, err := r.ReadString('\n')
		if err != nil {
			return nil, newConnError("Failed to get Error string in call to ReadString: %s", err)
		}
		return nil, parseError(cleanEnding(s))
	case "+":
		tmp, err := r.ReadString('\n')
		return cleanEnding(tmp), err
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
		return nil, newConnError("Failed to get raw int in call to ReadString: %s", err)
	}
	i, err := toInt(cleanEnding(s))
	if err != nil {
		return nil, newConversionError("Failed to convert raw int to int: %s", err)
	}
	return i, nil
}

func decodeBulkStringSuffix(r *bufio.Reader) (interface{}, error) {
	tmp, err := r.ReadString('\n')
	if err != nil {
		return nil, newConnError("Failed to get raw int for Bulk String size in call to ReadString: %s", err)
	}
	slen, err := toInt(cleanEnding(tmp))
	if err != nil {
		return nil, newConversionError("Failed to convert raw int to int for Bulk String size: %s", err)
	}
	if slen == nilLength {
		return RedisNil{}, nil
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
	tmp, err := r.ReadString('\n')
	if err != nil {
		return nil, newConnError("Failed to get raw int for Bulk Array size in call to ReadString: %s", err)
	}
	alen, err := toInt(cleanEnding(tmp))
	if err != nil {
		return nil, newConversionError("Failed to convert raw int to in for Bulk Array size: %s", err)
	}
	if alen == nilLength {
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

func toInt(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}
