package bot

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"go-transaction-bot/config"
	"strconv"
	"sync"
)

type Confirmation struct {
	SellerID      string
	BuyerID       string
	SellerPressed bool
	BuyerPressed  bool
	mu            *sync.Mutex
}

var (
	transConf = make(map[string]*Confirmation)

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"seller": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: func() string {
						transConf[i.ChannelID].SellerPressed = true
						transConf[i.ChannelID].SellerID = i.Member.User.ID
						return "User " + i.Member.Mention() + " is now a Seller!"
					}(),
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{
							discordgo.AllowedMentionTypeEveryone,
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"buyer": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: func() string {
						transConf[i.ChannelID].BuyerPressed = true
						transConf[i.ChannelID].BuyerID = i.Member.User.ID
						return "User " + i.Member.Mention() + " is now a Buyer!"
					}(),
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{
							discordgo.AllowedMentionTypeEveryone,
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"sell": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Selling some stuff?",
					// Buttons and other components are specified in Components field.
					Components: []discordgo.MessageComponent{
						// ActionRow is a container of all buttons within the same row.
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									// Label is what the user will see on the button.
									Label: "Confirm transaction as Seller",
									// Style provides coloring of the button. There are not so many styles tho.
									Style: discordgo.SuccessButton,
									// Disabled allows bot to disable some buttons for users.
									Disabled: false,
									// CustomID is a thing telling Discord which data to send when this button will be pressed.
									CustomID: "seller",
								},
								discordgo.Button{
									Label:    "Confirm transaction as Buyer",
									Style:    discordgo.PrimaryButton,
									Disabled: false,
									CustomID: "buyer",
								},
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
		},
		"confirm": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ch, err := s.Channel(i.ChannelID)
			if err != nil {
				fmt.Println("Can't return channel")
				return
			}

			if transConf[i.ChannelID].SellerPressed == true && transConf[i.ChannelID].BuyerPressed == true {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: func() string {
							q := `UPDATE sellbuyinfo SET selled = selled + 1 WHERE userid = $1;`
							Psql.pool.Exec(context.Background(), q, transConf[i.ChannelID].SellerID)
							q = `UPDATE sellbuyinfo SET bought = bought + 1 WHERE userid = $1;`
							Psql.pool.Exec(context.Background(), q, transConf[i.ChannelID].BuyerID)
							return "Apllying info to db and archiving the thread..."
						}(),
					},
				})
				if err != nil {
					panic(err)
				}
				s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
					Name:                 "[Sold] " + ch.Name,
					Topic:                ch.Topic,
					NSFW:                 ch.NSFW,
					Position:             ch.Position,
					Bitrate:              ch.Bitrate,
					UserLimit:            ch.UserLimit,
					PermissionOverwrites: ch.PermissionOverwrites,
					ParentID:             ch.ParentID,
					RateLimitPerUser:     ch.RateLimitPerUser,
					Archived:             true,
					AutoArchiveDuration:  60,
					Locked:               false,
					Invitable:            false,
				})
			} else {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Seller and Buyer must click on their buttons!",
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
		"Stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
					Content: func() string {
						res := struct {
							sold   int
							bought int
						}{}
						q := `SELECT selled, bought FROM sellbuyinfo WHERE usernickname = $1`
						row := Psql.pool.QueryRow(context.Background(), q, i.Member.User.Username+"#"+i.Member.User.Discriminator)
						if row.Scan(&res.sold, &res.bought) != nil {
							return "User has no stats!"
						} else {
							return "Bouhgt: " + strconv.Itoa(res.bought) + "\nSold: " + strconv.Itoa(res.sold)
						}
					}(),
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func handlers(goBot *discordgo.Session, botIsUp chan struct{}) {
	//start handler
	goBot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Info("Bot is up!")
		close(botIsUp)
	})
	//memberadd
	goBot.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		//addmem to database
	})
	//thread
	//goBot.AddHandler(func(s *discordgo.Session, t *discordgo.ThreadCreate) {
	//	_, _ = s.ChannelMessageSend(t.ID, "New thread!")
	//})

	//pingpong
	goBot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == config.Config.Discord.BotID {
			return
		}
		if m.Content == "ping" {
			_, _ = s.ChannelMessageSend(m.ChannelID, m.Author.Username)
			_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
		}
	})
	//buttons and commands
	goBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println("HERE")
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				fmt.Println("Command invoked")
				if _, inMap := transConf[i.ChannelID]; !inMap {
					transConf[i.ChannelID] = &Confirmation{
						SellerID:      "",
						SellerPressed: false,
						BuyerID:       "",
						BuyerPressed:  false,
						mu:            &sync.Mutex{},
					}
				}
				h(s, i)
				fmt.Println("chan struct = ", transConf[i.ChannelID])
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				fmt.Println("Button invoked")
				h(s, i)
				fmt.Println("chan struct after press= ", transConf[i.ChannelID])
			}
		}
	})

}
