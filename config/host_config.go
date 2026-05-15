package config

import (
	"errors"
	models "nxt_match_event_manager_api/internal/models/config"
)

/**
* Custom Host Configuration & Enviroment Setup
* @Maintainer : smeag_dev
* @Version : 1.0.0
* @Description: Define Setup Configuration {Host, Host, ENV}
**/
func GetHostConfig(env string) (*models.HostConfig, error) {
	switch env {
	case "dev":
		return models.NewHostConfig(
			Envs.PublicHost, // App Host
			Envs.Port,
			env,
		), nil

	case "stage":
		return models.NewHostConfig(
			Envs.AppStagHost,
			Envs.Port,
			env,
		), nil

	default:
		return nil, errors.New("unknown environment: " + env)
	}
}
