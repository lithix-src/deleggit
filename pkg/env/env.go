package env

import (
	"os"
	"strconv"
)

// Get returns the value of the environment variable named by the key.
// If the variable is not present, it returns the defaultValue.
func Get(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetInt returns the value as an integer. Returns default if missing or invalid.
func GetInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
