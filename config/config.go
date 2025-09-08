package config
import(
	"fmt"
	"os"
)
type Config struct{
	DBUser string
	DBPassword string
	DBHost string
	DBPort string
	DBName string
}
func LoadConfig() (*Config, error){
	cfg := &Config{
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBName: getEnv("DB_NAME", "demo"),
	}
	if cfg.DBPassword == ""{
		return nil, fmt.Errorf("DB Password must be set");
	}
	return cfg, nil;
}

func getEnv(key, fallback string) string{
	if value, ok := os.LookupEnv(key); ok{
		return value;
	}
	return fallback;
}