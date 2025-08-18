package userController

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-otp-auth-service/initializers"
	"math/big"
	"os"
	"strconv"
	"time"
)

func GenerateOTP(phone string) (string, error) {
	otpExpMin, err := strconv.Atoi(os.Getenv("OTP_EXP_MIN"))
	if err != nil {
		otpExpMin = 1
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	otp := fmt.Sprintf("%06d", n.Int64())
	err = initializers.RDB.Set(initializers.Ctx, "otp:"+phone, otp, time.Duration(otpExpMin)*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return otp, nil
}

func VerifyOTP(phone string, otp string) bool {
	storedOtp, err := initializers.RDB.Get(initializers.Ctx, "otp:"+phone).Result()
	if err != nil {
		return false
	}
	if storedOtp != otp {
		return false
	}
	initializers.RDB.Del(initializers.Ctx, "otp:"+phone)
	return true
}

func CheckOTPRequest(phone string, limit int, period time.Duration) (bool, string, error) {
	key := "otp_limit:" + phone
	count, err := initializers.RDB.Get(initializers.Ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, "", err
	}
	if count >= limit {
		ttl, _ := initializers.RDB.TTL(initializers.Ctx, key).Result()
		return false, fmt.Sprintf("please wait %s", ttl.Truncate(time.Second)), nil
	}
	pipe := initializers.RDB.TxPipeline()
	pipe.Incr(initializers.Ctx, key)
	if count == 0 {
		pipe.Expire(initializers.Ctx, key, period)
	}
	_, err = pipe.Exec(initializers.Ctx)
	if err != nil {
		return false, "", err
	}
	return true, "", nil
}
