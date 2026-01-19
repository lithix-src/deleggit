package env

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	key := "TEST_ENV_KEY"
	val := "test_val"

	os.Setenv(key, val)
	defer os.Unsetenv(key)

	if got := Get(key, "default"); got != val {
		t.Errorf("Get() = %v, want %v", got, val)
	}

	if got := Get("NON_EXISTENT_KEY", "default"); got != "default" {
		t.Errorf("Get() default = %v, want %v", got, "default")
	}
}

func TestGetInt(t *testing.T) {
	key := "TEST_ENV_INT"
	val := "123"

	os.Setenv(key, val)
	defer os.Unsetenv(key)

	if got := GetInt(key, 0); got != 123 {
		t.Errorf("GetInt() = %v, want %v", got, 123)
	}

	if got := GetInt("NON_EXISTENT_KEY", 42); got != 42 {
		t.Errorf("GetInt() default = %v, want %v", got, 42)
	}

	os.Setenv(key, "not_an_int")
	if got := GetInt(key, 99); got != 99 {
		t.Errorf("GetInt() invalid = %v, want %v", got, 99)
	}
}
