package config

import (
	"errors"
	models "nxt_match_event_manager_api/internal/models/config"
)

/**
* Custom Database Configuration & Enviroment Setup
* @Maintainer : nxt_match_dev
* @Version : 1.0.0
* @Description: Define Setup Configuration {Host, URI, DBName, Host, Port}
**/
func GetDBConfig(env string) (*models.DBconfig, error) {
	switch env {
	case "dev":
		return models.NewDBConfig(
			Envs.MongoLocalUri,
			Envs.DBName,
			Envs.DBHost, // DB host
			Envs.DBPort,
			env,
		), nil

	case "stage":
		return models.NewDBConfig(
			Envs.MongoStagUri,
			Envs.MongoStagName,
			Envs.MongoStagIpAddr,
			Envs.MongoStagPort,
			env,
		), nil

	default:
		return nil, errors.New("unknown environment: " + env)
	}
}
