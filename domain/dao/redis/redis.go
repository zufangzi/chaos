package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"opensource/chaos/common"
	"reflect"
	"runtime/debug"
)

var c *redis.Conn

var MAX_POOL_SIZE = 20
var redisPoll chan redis.Conn
var redisServer string

func RedisInit() {
	redisServer = common.Path.RedisUrl
}

func putRedis(conn redis.Conn) {
	if redisPoll == nil {
		redisPoll = make(chan redis.Conn, MAX_POOL_SIZE)
	}
	if len(redisPoll) >= MAX_POOL_SIZE {
		conn.Close()
		return
	}
	redisPoll <- conn
}

func DftGetRedisConn() redis.Conn {
	return GetRedisConn("tcp", redisServer)
}

func GetRedisConn(network, address string) redis.Conn {
	if len(redisPoll) == 0 {
		redisPoll = make(chan redis.Conn, MAX_POOL_SIZE)
		go func() {
			for i := 0; i < MAX_POOL_SIZE/2; i++ {
				c, err := redis.Dial(network, address)
				common.AssertPanic(err)
				putRedis(c)
			}
		}()
	}
	// 可能会造成阻塞？
	return <-redisPoll
}

func Safe(f func(redis.Conn)) {
	c := DftGetRedisConn()
	defer func() {
		putRedis(c)
		if e, ok := recover().(error); ok {
			log.Println("catchable redis error occur. " + e.Error())
			debug.PrintStack()
		}
	}()
	f(c)
}

func simpleMultiArgsCmd(c redis.Conn, cmdStr string, container string, args ...string) []reflect.Value {
	t := reflect.ValueOf(c)
	m := t.MethodByName("Do")
	var values []reflect.Value
	if container != "" {
		values = make([]reflect.Value, len(args)+2)
	} else {
		values = make([]reflect.Value, len(args)+1)
	}

	values[0] = reflect.ValueOf(cmdStr)
	i := 1
	if container != "" {
		values[i] = reflect.ValueOf(container)
		i++
	}
	for _, v := range args {
		values[i] = reflect.ValueOf(v)
		i++
	}
	return m.Call(values)
}

func Lpush(queue string, value ...string) bool {
	result := false
	Safe(func(c redis.Conn) {
		rst := simpleMultiArgsCmd(c, "LPUSH", queue, value...)
		ok, err := redis.Bool(rst[0].Interface().(int64), nil)
		common.AssertPrint(err)
		result = ok
	})
	return result
}

func Rpop(queue string) string {
	var result string
	Safe(func(c redis.Conn) {
		rst := simpleMultiArgsCmd(c, "RPOP", queue)
		data, err := redis.String(rst[0].Interface(), nil)
		common.AssertPrint(err)
		result = data
	})
	return result
}
