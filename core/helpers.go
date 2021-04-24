package core

import "os"

/* Nice little function to get an env var with fallback */
func GetEnv(name string, fallback string) string {
	value, present := os.LookupEnv(name)
	if present {
		return value
	}
	return fallback
}
