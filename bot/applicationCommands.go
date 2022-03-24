package bot

import (
	"github.com/bwmarrin/discordgo"
	"go-transaction-bot/config"
	"log"
)

func createCommands(goBot *discordgo.Session) {
	createSellCommand(goBot)
	createConfirmCommand(goBot)
	createStatCommand(goBot)
}

func createStatCommand(goBot *discordgo.Session) {
	stat, err := goBot.ApplicationCommandCreate(config.Config.Discord.AppID, config.Config.Discord.GuildID, &discordgo.ApplicationCommand{
		Name: "Stats",
		Type: discordgo.UserApplicationCommand,
	})
	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	err = goBot.ApplicationCommandPermissionsEdit(config.Config.Discord.AppID, config.Config.Discord.GuildID, stat.ID, &discordgo.ApplicationCommandPermissionsList{
		Permissions: []*discordgo.ApplicationCommandPermissions{
			&discordgo.ApplicationCommandPermissions{
				ID:         config.Config.Discord.MemberRoleID,
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
			&discordgo.ApplicationCommandPermissions{
				ID:         config.Config.Discord.ModeratorRoleID,
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Cannot modify sell perm command: %v", err)
	}
}

func createSellCommand(goBot *discordgo.Session) {
	sell, err := goBot.ApplicationCommandCreate(config.Config.Discord.AppID, config.Config.Discord.GuildID, &discordgo.ApplicationCommand{
		Name:        "sell",
		Description: "Starting trading process",
	})
	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	err = goBot.ApplicationCommandPermissionsEdit(config.Config.Discord.AppID, config.Config.Discord.GuildID, sell.ID, &discordgo.ApplicationCommandPermissionsList{
		Permissions: []*discordgo.ApplicationCommandPermissions{
			&discordgo.ApplicationCommandPermissions{
				ID:         config.Config.Discord.ModeratorRoleID,
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
			&discordgo.ApplicationCommandPermissions{
				ID:         config.Config.Discord.MemberRoleID,
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: false,
			},

			&discordgo.ApplicationCommandPermissions{
				ID:         "953702196476792832",
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
			&discordgo.ApplicationCommandPermissions{
				ID:         "950582530854228019",
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Cannot modify sell perm command: %v", err)
	}
}

func createConfirmCommand(goBot *discordgo.Session) {
	confirm, err := goBot.ApplicationCommandCreate(config.Config.Discord.AppID, config.Config.Discord.GuildID, &discordgo.ApplicationCommand{
		Name:        "confirm",
		Description: "Transaction confirmation",
	})
	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	err = goBot.ApplicationCommandPermissionsEdit(config.Config.Discord.AppID, config.Config.Discord.GuildID, confirm.ID, &discordgo.ApplicationCommandPermissionsList{
		Permissions: []*discordgo.ApplicationCommandPermissions{
			&discordgo.ApplicationCommandPermissions{
				ID:         config.Config.Discord.ModeratorRoleID,
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
			&discordgo.ApplicationCommandPermissions{
				ID:         config.Config.Discord.MemberRoleID,
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: false,
			},
			&discordgo.ApplicationCommandPermissions{
				ID:         "953702196476792832",
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
			&discordgo.ApplicationCommandPermissions{
				ID:         "950582530854228019",
				Type:       discordgo.ApplicationCommandPermissionTypeRole,
				Permission: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Cannot modify confirm perm command: %v", err)
	}
}
