package internal

import "os"

func GetServerAddr() string {
	addr, ok := os.LookupEnv("HTTP_ADDR")
	if !ok {
		return "localhost:8080"
	}
	return addr
}
