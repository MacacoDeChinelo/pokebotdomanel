package twitch

import (
	"context"
	"darthverde/database"
	"darthverde/models"
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "settwitch",
			Description: "Configura o alerta de live da Twitch para este servidor",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "canal",
					Description: "Canal onde o aviso será enviado",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "cargo",
					Description: "Cargo que será mencionado (Ex: @Inscritos)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "streamer",
					Description: "Nome do canal na Twitch (ex: gaules)",
					Required:    true,
				},
			},
		},
	}
}

func HandleSetTwitch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Montando os dados para salvar
	alerta := models.TwitchAlert{
		ServerID:       i.GuildID,
		DiscordChannel: optionMap["canal"].ChannelValue(s).ID,
		RoleToMention:  optionMap["cargo"].RoleValue(s, i.GuildID).ID,
		TwitchChannel:  optionMap["streamer"].StringValue(),
		IsLive:         false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.InsertTwitchAlert(ctx, alerta)

	msg := "✅ Alerta da Twitch configurado com sucesso para o canal: " + alerta.TwitchChannel
	if err != nil {
		msg = "❌ Erro ao salvar o alerta da Twitch no banco de dados."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
