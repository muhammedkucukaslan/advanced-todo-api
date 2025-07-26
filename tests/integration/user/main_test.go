package testuser

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	fmt.Println("Running user integration tests...")
	err := godotenv.Load("../../../.env")
	if err != nil {
		panic("Failed to load .env file: " + err.Error())
	}

	code := m.Run()

	os.Exit(code)
}
