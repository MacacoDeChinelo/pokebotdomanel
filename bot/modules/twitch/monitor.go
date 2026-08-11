package twitch

import (
	"bytes"
	"context"
	"darthverde/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	libdatabase "github.com/MacacoDeChinelo/readconf/pkg/database"
	"github.com/bwmarrin/discordgo"
)

// Estruturas para ler a resposta da Twitch
type TwitchAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type TwitchStreamResponse struct {
	Data []struct {
		Type string `json:"type"` // "live" indica que tá online
	} `json:"data"`
}

// COLOQUE SUAS CREDENCIAIS AQUI
var (
	twitchClientID,
	twitchClientSecret string
)
var accessToken string

func StartTwitchMonitor(s *discordgo.Session) {
	// COLOQUE SUAS CREDENCIAIS AQUI

	twitchClientID = libdatabase.GetVariable("twitchClient").(string)
	twitchClientSecret = libdatabase.GetVariable("twitchSecret").(string)

	fmt.Println(time.Now().Format("02-01-2006 15:04:05") + " Iniciando monitoramento de lives da Twitch...")

	// Roda a cada 5 minutos
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for range ticker.C {
			checkTwitchLives(s)
		}
	}()
}

// Essa função pega o token obrigatório para consultar a Twitch
func getTwitchToken() (string, error) {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", twitchClientID, twitchClientSecret)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var auth TwitchAuthResponse
	json.NewDecoder(resp.Body).Decode(&auth)
	return auth.AccessToken, nil
}

func checkTwitchLives(s *discordgo.Session) {
	// Se não temos o token da Twitch, vamos gerar um
	if accessToken == "" {
		token, err := getTwitchToken()
		if err != nil {
			log.Println("Erro ao gerar token da Twitch:", err)
			return
		}
		accessToken = token
	}

	ctx := context.Background()
	alertas, err := database.GetAllTwitchAlerts(ctx)
	if err != nil {
		log.Println("Erro ao buscar alertas da Twitch no banco:", err)
		return
	}

	for _, alerta := range alertas {
		// Monta a URL de consulta para aquele streamer específico
		url := "https://api.twitch.tv/helix/streams?user_login=" + alerta.TwitchChannel

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Client-ID", twitchClientID)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)

		// Se der erro 401 (Não autorizado), o token venceu. Limpamos para gerar de novo no próximo ciclo.
		if err != nil || resp.StatusCode == 401 {
			accessToken = ""
			continue
		}

		var streamData TwitchStreamResponse
		json.NewDecoder(resp.Body).Decode(&streamData)
		resp.Body.Close()

		isLiveNow := len(streamData.Data) > 0 && streamData.Data[0].Type == "live"

		// Se está online agora mas no banco estava offline -> Envia mensagem!
		if isLiveNow && !alerta.IsLive {
			mensagem := fmt.Sprintf("Fala galera <@&%s>, o canal **%s** acabou de entrar em live!\nCorre pra assistir: https://twitch.tv/%s", alerta.RoleToMention, alerta.TwitchChannel, alerta.TwitchChannel)
			s.ChannelMessageSend(alerta.DiscordChannel, mensagem)

			// Atualiza no banco que já foi avisado
			database.UpdateTwitchLiveStatus(ctx, alerta.ID, true)
		}

		// Se está offline agora mas no banco estava online -> Atualiza para offline
		if !isLiveNow && alerta.IsLive {
			database.UpdateTwitchLiveStatus(ctx, alerta.ID, false)
		}
	}
}
