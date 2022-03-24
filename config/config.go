package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Config *config = &config{}
)

type config struct {
	Discord struct {
		Token           string `json : "Token"`
		BotPrefix       string `json : "BotPrefix"`
		GuildID         string `json : "GuildID"`
		AppID           string `json : "AppID"`
		ModeratorRoleID string `json : "ModeratorRoleID"`
		MemberRoleID    string `json : "ModeratorRoleID"`
		BotID           string
	} `json : "Discord"`
	Postgres struct {
		Username string `json : "Username"`
		Password string `json : "Password"`
		Host     string `json : "Host"`
		Port     string `json : "Port"`
		Database string `json : "Database"`
	}
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(file))
	err = json.Unmarshal(file, Config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
