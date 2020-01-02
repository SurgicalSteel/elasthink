package util

import (
	"strings"
)

//GetEnv gets environment from given string
func GetEnv(env string) string {
	env = strings.ToLower(env)
	switch env {
	case "stg", "staging":
		return "staging"
	case "prod", "production":
		return "production"
	default:
		return "development"
	}
}
