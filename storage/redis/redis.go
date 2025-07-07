package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func GetOtpFromMemoryOrRedis(email string) string {
	key := "otp:" + email
	otp, err := Client.Get(Ctx, key).Result()
	if err != nil {
		return "" // OTP tidak ditemukan
	}
	return otp
}

func StoreOtpInMemoryOrRedis(email, otp string) {
	key := "otp:" + email
	err := Client.Set(Ctx, key, otp, 0).Err()
	if err != nil {
		panic(err) // Gagal menyimpan OTP ke Redis
	}
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
