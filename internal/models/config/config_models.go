package models

type DBconfig struct {
	Uri      string
	DBName   string
	DBHost   string
	DBPort   string
	BuildEnv string
}

func NewDBConfig(uri string, dbName string, DbHost string, DbPort string, BuildEnv string) *DBconfig {

	return &DBconfig{
		Uri:      uri,
		DBName:   dbName,
		DBHost:   DbHost,
		DBPort:   DbPort,
		BuildEnv: BuildEnv,
	}
}
