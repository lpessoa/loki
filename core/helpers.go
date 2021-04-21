package core

import "os"

func GetEnv(name string, fallback string) string {
	value, present := os.LookupEnv(name)
	if present {
		return value
	}
	return fallback
}
