package youtube

import (
	"context"
	"darthverde/database"
	"darthverde/models"
	"fmt"

	"github.com/bwmarrin/discordgo"
	// Lembre-se de importar o seu pacote models e database aqui
	// "seu_projeto/models"
	// "seu_projeto/database"
)

// HandleSetYouTube processa o comando /setyoutube
func HandleSetYouTube(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 1. Extrair as opções enviadas pelo usuário
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// 2. Pegar os valores específicos
	canalDiscord := optionMap["canal_texto"].ChannelValue(s)
	cargoMencao := optionMap["cargo"].RoleValue(s, i.GuildID)
	idYouTube := optionMap["id_youtube"].StringValue()

	// 3. Montar o modelo de dados que definimos anteriormente
	alerta := models.YouTubeAlert{
		ServerID:       i.GuildID,
		DiscordChannel: canalDiscord.ID,
		RoleToMention:  cargoMencao.ID,
		YouTubeChannel: idYouTube,
		IsLive:         false, // Começa como false por padrão
	}

	// 4. Salvar no MongoDB
	// Aqui você chama a função do seu pacote de banco de dados
	err := database.InsertYouTubeAlert(context.Background(), alerta)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Ocorreu um erro ao salvar no banco de dados.",
				Flags:   discordgo.MessageFlagsEphemeral, // Mensagem invisível para os outros
			},
		})
		return
	}

	// 5. Responder com sucesso usando um Embed
	embed := &discordgo.MessageEmbed{
		Title:       "🔴 Alerta de Live Configurado!",
		Color:       0xFF0000, // Vermelho do YouTube
		Description: "O bot agora vai avisar automaticamente quando o canal entrar em live.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Canal de Aviso",
				Value:  fmt.Sprintf("<#%s>", canalDiscord.ID),
				Inline: true,
			},
			{
				Name:   "Cargo Mencionado",
				Value:  fmt.Sprintf("<@&%s>", cargoMencao.ID),
				Inline: true,
			},
			{
				Name:   "ID do YouTube",
				Value:  fmt.Sprintf("`%s`", idYouTube),
				Inline: false,
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
