package config

import (
	"encoding/json"
	"os"
)

type JWTConfig struct {
	SecretKey string `json:"secret_key"`
}

func GetJWTConfig() JWTConfig {
	jwtConfig := JWTConfig{}
	jsonFile, err := os.Open("././configs/JwtConfig.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&jwtConfig)
	if err != nil {
		panic(err)
	}
	return jwtConfig
}
