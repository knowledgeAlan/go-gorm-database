package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"vscode/go-gorm-database/models"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var RedisClient *redis.Client

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis init err", err)
	}

	fmt.Println("ping", pong)
}

func SetKey(key string, value string) {
	err := RedisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println("set key", err)
	}
}

func GetKey(key string) string {
	res, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("get key error", err)
	}
	return res
}

func CreateAuth(userId int64, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)

	now := time.Now()

	errAccess := RedisClient.Set(ctx, td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()

	if errAccess != nil {
		return errAccess
	}

	errRefresh := RedisClient.Set(ctx, td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()

	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

//获取登录信息
func FetchAuth(authD *models.AccessDetails) (int64, error) {

	userId, err := RedisClient.Get(ctx, authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userId, 10, 64)

	return userID, nil
}

func DeleteAuth(giveUuid string) (int64, error) {

	deleted, err := RedisClient.Del(ctx, giveUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
