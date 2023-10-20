package shared

import (
	"os"
	"strconv"
)

func MustGetEnv(key string) string {
	if v := os.Getenv(key); v == "" {
		panic("Missing required environment variable: " + key)
	} else {
		return v
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func StringToUint(s string) (uint, error) {
	// Convert string to uint64
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	// Cast uint64 to uint and return
	return uint(u64), nil
}
