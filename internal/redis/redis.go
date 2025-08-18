package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(addr string, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis ✔️")
}

func Set(key string, value string, expiration time.Duration) error {
	return Rdb.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return Rdb.Get(Ctx, key).Result()
}

// Bitfield operations for checkbox storage
const CheckboxBitfield = "checkboxes"

func SetCheckbox(index int, value bool) error {
	bitValue := 0
	if value {
		bitValue = 1
	}
	return Rdb.BitField(Ctx, CheckboxBitfield, "SET", "u1", index, bitValue).Err()
}

func GetCheckbox(index int) (bool, error) {
	result := Rdb.BitField(Ctx, CheckboxBitfield, "GET", "u1", index)
	values, err := result.Result()
	if err != nil {
		return false, err
	}
	if len(values) == 0 {
		return false, nil // Default to unchecked
	}
	return values[0] == 1, nil
}

func GetAllCheckboxes(maxIndex int) ([]bool, error) {
	checkboxes := make([]bool, maxIndex)
	for i := 0; i < maxIndex; i++ {
		checked, err := GetCheckbox(i)
		if err != nil {
			return nil, err
		}
		checkboxes[i] = checked
	}
	return checkboxes, nil
}
