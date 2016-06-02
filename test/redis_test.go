package test

import (
	"opensource/chaos/background/server/domain/redis"
	"testing"
)

func Test_redis_lpush(t *testing.T) {
	if !redis.Lpush("hello_new5", "world1", "world2") {
		t.Fail()
	}
}
