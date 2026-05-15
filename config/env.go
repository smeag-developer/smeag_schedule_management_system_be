package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost           string
	Port                 string
	DBUser               string
	DBPort               string
	DBHost               string
	DBPassword           string
	DBAddress            string
	DBName               string
	MongoLocalUri        string
	MongoStagUri         string
	MongoStagName        string
	MongoStagIpAddr      string
	AppStagHost          string
	MongoStagPort        string
	Dev_Allowed_Origins  string
	Stag_Allowed_Origins string
}

/*Create Singleton*/
var Envs = initConfig()

func initConfig() Config {

	//reload enviroment var
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", os.Getenv("PUBLIC_HOST")),
		Port:       getEnv("PORT", os.Getenv("PORT")),
		DBUser:     getEnv("MONGO_LOCAL_DB_USER", os.Getenv("MONGO_LOCAL_DB_USER")),
		DBPort:     getEnv("MONGO_LOCAL_DB_PORT", os.Getenv("MONGO_LOCAL_DB_PORT")),
		DBHost:     getEnv("MONGO_LOCAL_DB_HOST", os.Getenv("MONGO_LOCAL_DB_HOST")),
		DBPassword: getEnv("MONGO_LOCAL_DB_PASSWORD", os.Getenv("MONGO_LOCAL_DB_PASSWORD")),
		DBAddress: fmt.Sprintf("%s:%s",
			getEnv("MONGO_LOCAL_DB_HOST", os.Getenv("MONGO_LOCAL_DB_HOST")),
			getEnv("MONGO_LOCAL_DB_PORT", os.Getenv("MONGO_LOCAL_DB_PORT"))),
		DBName:               getEnv("MONGO_LOCAL_DB_NAME", os.Getenv("MONGO_LOCAL_DB_NAME")),
		MongoLocalUri:        getEnv("MONGO_LOCAL_URI", os.Getenv("MONGO_LOCAL_URI")),
		MongoStagUri:         getEnv("MONGO_DB_STAG_URI", os.Getenv("MONGO_DB_STAG_URI")),
		MongoStagIpAddr:      getEnv("MONGO_DB_STAG_IP_ADDR", os.Getenv("MONGO_DB_STAG_IP_ADDR")),
		MongoStagPort:        getEnv("MONGO_DB_STAG_PORT", os.Getenv("MONGO_DB_STAG_PORT")),
		AppStagHost:          getEnv("APP_STAG_HOST", os.Getenv("APP_STAG_HOST")),
		MongoStagName:        getEnv("MONGO_DB_STAG_NAME", os.Getenv("MONGO_DB_STAG_NAME")),
		Dev_Allowed_Origins:  getEnv("DEV_ALLOWED_ORIGINS", os.Getenv("DEV_ALLOWED_ORIGINS")),
		Stag_Allowed_Origins: getEnv("STAG_ALLOWED_ORIGINS", os.Getenv("STAG_ALLOWED_ORIGINS")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
