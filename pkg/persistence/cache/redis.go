package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"time"
)

const (
	lockExpireTime = 10 * time.Second
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(dsn string) (*RedisCache, error) {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)
	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}, nil
}

func (c *RedisCache) TryLock(key string) error {
	logrus.Debugf("try lock key: %s", key)
	if ok, _ := c.client.SetNX(c.ctx, key, true, lockExpireTime).Result(); !ok {
		logrus.Infof("key %s is locked", key)
		return errors.New("already locked")
	}
	return nil
}

func (c *RedisCache) Unlock(key string) error {
	if _, err := c.client.Del(c.ctx, key).Result(); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Subscribe(topic string, fn interface{}) error {
	logrus.Debugf("subscribe topic: %s", topic)
	subject := c.client.Subscribe(c.ctx, topic)
	go func() {
		defer subject.Close()

		for {
			msg, err := subject.ReceiveMessage(c.ctx)
			if err != nil {
				logrus.Errorf("receive message error: %v", err)
			}

			// 反射执行函数
			function := reflect.ValueOf(fn)
			args := function.Type()
			var params []reflect.Value
			for i := 0; i < args.NumIn(); i++ {
				param := reflect.New(args.In(i))
				switch args.In(i).Kind() {
				case reflect.String:
					param.Elem().Set(reflect.ValueOf(msg.Payload))
				case reflect.Uint:
					parseUint, err := strconv.ParseUint(msg.Payload, 10, 64)
					if err != nil {
						logrus.Errorf("parse uint error when subscribe %s: %v", topic, err)
					}
					param.Elem().Set(reflect.ValueOf(uint(parseUint)))
				}
				params = append(params, param.Elem())
			}

			function.Call(params)
		}
	}()
	return nil
}

func (c *RedisCache) Publish(topic string, data interface{}) {
	c.client.Publish(c.ctx, topic, data)
}
