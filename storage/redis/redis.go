package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "redis:localhost:6379"
	}
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		fmt.Println("❌ Failed to parse Redis URL:", err)
		return
	}
	Client = redis.NewClient(opt)
}

func GetOtpFromMemoryOrRedis(email string) string {
	key := "otp:" + email
	otp, err := Client.Get(Ctx, key).Result()
	if err != nil {
		return "" // OTP tidak ditemukan
	}
	return otp
}

func StoreOtpInMemoryOrRedis(email, otp string) error {
	key := "otp:" + email
	err := Client.Set(Ctx, key, otp, 0).Err()
	if err != nil {
		fmt.Println("❌ Failed storing OTP:", err)
		return err
	}
	return nil
}

func BlacklistToken(token string, ttl time.Duration) error {
	return Client.Set(Ctx, "blacklist:"+token, "1", ttl).Err()
}

func IsTokenBlacklisted(token string) (bool, error) {
	val, err := Client.Get(Ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
