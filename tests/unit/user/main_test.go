package testuser

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {

	err := godotenv.Load("../../../.env")
	if err != nil {
		panic("Failed to load .env file: " + err.Error())
	}

	code := m.Run()

	os.Exit(code)
}
