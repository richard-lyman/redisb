package redisb

import (
	"fmt"
	"io"
)

func prepend(s string, ss []string) []string {
	return append([]string{s}, ss...)
}

// int - incr decr incrby decrby strlen
func Incr(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("incr", args)...)
}
func Decr(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("decr", args)...)
}
func Incrby(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("incrby", args)...)
}
func Decrby(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("decrby", args)...)
}

// bool - msetnx
func Msetnx(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("msetnx", args)...)
}

// string - incrbyfloat mset
func Incrbyfloat(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("incrbyfloat", args)...)
}
func Mset(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("mset", args)...)
}

// array - mget
func Mget(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("mget", args)...)
}

// LIST
// int - LINSERT LLEN LPUSH LPUSHX LREM RPUSH RPUSHX
func Linsert(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("linsert", args)...)
}
func Llen(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("llen", args)...)
}
func Lpush(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("lpush", args)...)
}
func Lpushx(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("lpushx", args)...)
}
func Lrem(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("lrem", args)...)
}
func Rpush(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("rpush", args)...)
}
func Rpushx(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("rpushx", args)...)
}

// bool - LSET LTRIM
func Lset(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("lset", args)...)
}
func Ltrim(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("ltrim", args)...)
}

// string - BRPOPLPUSH RPOPLPUSH LINDEX LPOP RPOP
func Brpoplpush(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("brpoplpush", args)...)
}
func Rpoplpush(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("rpoplpush", args)...)
}
func Lindex(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("lindex", args)...)
}
func Lpop(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("lpop", args)...)
}
func Rpop(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("rpop", args)...)
}

// array - BLPOP BRPOP LRANGE
func Blpop(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("blpop", args)...)
}
func Brpop(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("brpop", args)...)
}
func Lrange(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("lrange", args)...)
}

// SET
// int - SADD SCARD SDIFFSTORE SINTERSTORE SREM SUNIONSTORE
func Sadd(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("sadd", args)...)
}
func Scard(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("scard", args)...)
}
func Sdiffstore(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("sdiffstore", args)...)
}
func Sinterstore(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("sinterstore", args)...)
}
func Srem(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("srem", args)...)
}
func Sunionstore(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("sunionstore", args)...)
}

// bool - SISMEMBER SMOVE
func Sismember(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("sismember", args)...)
}
func Smove(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("smove", args)...)
}

// string - SPOP
func Spop(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("spop", args)...)
}

// array - SDIFF SINTER SMEMBERS SUNION SSCAN
func Sdiff(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("sdiff", args)...)
}
func Sinter(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("sinter", args)...)
}
func Smembers(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("smembers", args)...)
}
func Sunion(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("sunion", args)...)
}
func Sscan(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("sscan", args)...)
}

// string or array - SRANDMEMBER
func Srandmember(rw io.ReadWriter, args ...string) (interface{}, error) {
	return Raw(rw, prepend("srandmember", args)...)
}

// SORTED SET
// string - ZADD ZINCRBY ZRANK ZREVRANK ZSCORE
func Zadd(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("zadd", args)...)
}
func Zincrby(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("zincrby", args)...)
}
func Zrank(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("zrank", args)...)
}
func Zrevrank(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("zrevrank", args)...)
}
func Zscore(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("zscore", args)...)
}

// int - ZCARD ZCOUNT ZINTERSTORE ZLEXCOUNT ZREM
//	ZREMRANGEBYLEX ZREMRANGEBYRANK ZREMRANGEBYSCORE ZUNIONSTORE
func Zcard(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zcard", args)...)
}
func Zcount(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zcount", args)...)
}
func Zinterstore(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zinterstore", args)...)
}
func Zlexcount(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zlexcount", args)...)
}
func Zrem(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zrem", args)...)
}
func Zremrangebylex(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zremrangebylex", args)...)
}
func Zremrangebyrank(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zremrangebyrank", args)...)
}
func Zremrangebyscore(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zremrangebyscore", args)...)
}
func Zunionstore(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("zunionstore", args)...)
}

// array - ZRANGE ZRANGEBYLEX ZREVRANGEBYLEX ZRANGEBYSCORE ZREVRANGE ZREVRANGEBYSCORE ZSCAN
func Zrange(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zrange", args)...)
}
func Zrangebylex(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zrangebylex", args)...)
}
func Zrevrangebylex(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zrevrangebylex", args)...)
}
func Zrevrangebyscore(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zrevrangebyscore", args)...)
}
func Zrevrange(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zrevrange", args)...)
}
func Zrangebyscore(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zrangebyscore", args)...)
}
func Zscan(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("zscan", args)...)
}

// HASH
// bool - HMSET HSET HSETNX
func Hmset(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("hmset", args)...)
}
func Hset(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("hset", args)...)
}
func Hsetnx(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("hsetnx", args)...)
}

// int - HDEL HEXISTS HINCRBY HLEN HSTRLEN
func Hdel(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("hdel", args)...)
}
func Hexists(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("hexists", args)...)
}
func Hincrby(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("hincrby", args)...)
}
func Hlen(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("hlen", args)...)
}
func Hstrlen(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("hstrlen", args)...)
}

// string - HGET HINCRBYFLOAT
func Hget(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("hget", args)...)
}
func Hincrbyfloat(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("hincrbyfloat", args)...)
}

// array - HGETALL HKEYS HMGET HVALS HSCAN
func Hgetall(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("hgetall", args)...)
}
func Hkeys(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("hkeys", args)...)
}
func Hmget(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("hmget", args)...)
}
func Hvals(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("hvals", args)...)
}
func Hscan(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("hscan", args)...)
}

// KEYS
// int DEL PTTL TOUCH TTL WAIT
func Del(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("del", args)...)
}
func Pttl(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("pttl", args)...)
}
func Touch(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("touch", args)...)
}
func Ttl(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("ttl", args)...)
}
func Wait(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("wait", args)...)
}

// bool EXISTS EXPIRE EXPIREAT MIGRATE MOVE PERSIST PEXPIRE PEXPIREAT RENAME RENAMENX RESTORE
func Exists(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("exists", args)...)
}
func Expire(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("expire", args)...)
}
func Expireat(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("expireat", args)...)
}
func Move(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("move", args)...)
}
func Persist(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("persist", args)...)
}
func Pexpire(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("pexpire", args)...)
}
func Pexpireat(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("pexpireat", args)...)
}
func Rename(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("rename", args)...)
}
func Renamenx(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("renamenx", args)...)
}

// string - RANDOMKEY TYPE
func Randomkey(rw io.ReadWriter) (string, error) { return String(rw, "randomkey") }
func RedisType(rw io.ReadWriter, args ...string) (string, error) {
	return String(rw, prepend("type", args)...)
}

// array KEYS SORT SCAN
func Keys(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("keys", args)...)
}
func Sort(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("sort", args)...)
}
func Scan(rw io.ReadWriter, args ...string) ([]interface{}, error) {
	return Array(rw, prepend("scan", args)...)
}

/// raw - OBJECT
func Object(rw io.ReadWriter, args ...string) (interface{}, error) {
	return Raw(rw, prepend("object", args)...)
}

// SCRIPTING
// EVAL EVALSHA
// SCRIPT DEBUG YES|SYNC|NO
// SCRIPT EXISTS
// SCRIPT FLUSH
// SCRIPT KILL
// SCRIPT LOAD
func Eval(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("eval", args)...)
}
func Evalsha(rw io.ReadWriter, args ...string) (int64, error) {
	return Int64(rw, prepend("evalsha", args)...)
}
func ScriptExists(rw io.ReadWriter, args ...string) ([]bool, error) {
	return Bools(rw, prepend("script", prepend("exists", args))...)
}
func ScriptLoad(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("script", prepend("load", args))...)
}

// TRANSACTION
// DISCARD EXEC MULTI UNWATCH WATCH
func Multi(rw io.ReadWriter) (io.ReadWriter, error) {
	result, err := Bool(rw, "multi")
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, fmt.Errorf("Failed to open transaction.")
	}
	return rw, nil
}

func Watch(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("watch", args)...)
}

func Unwatch(rw io.ReadWriter, args ...string) (bool, error) {
	return Bool(rw, prepend("unwatch", args)...)
}

func Exec(rw io.ReadWriter) ([]interface{}, error) {
	return Array(rw, "exec")
}

func Discard(rw io.ReadWriter) (bool, error) {
	return Bool(rw, "discard")
}
