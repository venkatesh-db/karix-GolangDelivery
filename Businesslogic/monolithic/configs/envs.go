package configs

import ( 
	   "os"
       "fmt"
	   "time"
)

func getenv(key, defaultVal string) string {

	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getenvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	var val int
	_, err := fmt.Sscanf(valStr, "%d", &val)
	if err != nil {
		return defaultVal
	}
	return val
}

func getenvDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
