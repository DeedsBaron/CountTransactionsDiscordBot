package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"go-transaction-bot/config"
)

var (
	Psql *Postgres
)

func Start() {
	Psql = NewPostgres()
	botIsUp := make(chan struct{})
	logrus.Info("Successfully connected to database")
	goBot, err := discordgo.New("Bot " + config.Config.Discord.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := goBot.User("@me")
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	config.Config.Discord.BotID = u.ID
	createCommands(goBot)
	handlers(goBot, botIsUp)
	err = goBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	<-botIsUp
	//AddmembersToDB(goBot, Psql)
}
