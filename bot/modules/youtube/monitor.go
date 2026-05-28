// bot/youtube/monitor.go
package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pokebot/database"
	"pokebot/utils"
	"time"

	"github.com/bwmarrin/discordgo"
	libdatabase "github.com/jolealpe89/readconf/pkg/database"
	// importe seus models e database
)

func StartYouTubeMonitor(s *discordgo.Session) {
	fmt.Println(time.Now().Format("02-01-2006 15:04:05") + " Iniciando monitoramento de lives do YouTube...")
	// Cria um Ticker para rodar a cada 5 minutos
	minutos := libdatabase.GetVariable("monitorTime").(int32)
	if minutos < 5 {
		minutos = 5 // valor padrão de 5 minutos se a variável não estiver configurada corretamente
		log.Printf("Valor inválido para monitorTime: %d. Usando valor padrão de 5 minutos.\n", minutos)
	}
	tempo := time.Duration(minutos) * time.Minute
	ticker := time.NewTicker(tempo)

	go func() {
		for range ticker.C {
			checkYouTubeLives(s)
		}
	}()
}

type YouTubeSearchResponse struct {
	Items []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

func checkYouTubeLives(s *discordgo.Session) {
	fmt.Println(time.Now().Format("02-01-2006 15:04:05") + " Verificando lives online...")
	// Pega a chave da API das variáveis de ambiente do seu sistema ou arquivo .env
	apiKeyEncrypt := libdatabase.GetVariable("youtubeApiKey").(string) //os.Getenv("YOUTUBE_API_KEY")
	apiKey, err := utils.Decrypt(apiKeyEncrypt)
	if apiKey == "" {
		log.Println("AVISO: Variável YOUTUBE_API_KEY não está configurada!")
		return
	}

	// 1. Buscar todos os alertas configurados no MongoDB que interagem com a collection streamer_alerts [cite: 28]
	alertas, err := database.GetAllYouTubeAlerts(context.Background())
	if err != nil {
		log.Printf("Erro ao buscar alertas do YouTube: %v\n", err)
		return
	}

	// 2. Iterar sobre cada configuração
	for _, alerta := range alertas {
		// 3. Montar a URL da API do YouTube utilizando o endpoint /search com os filtros corretos [cite: 14]
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&eventType=live&type=video&key=%s", alerta.YouTubeChannel, apiKey)

		// 4. Fazer a requisição HTTP GET
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Erro na requisição para o canal %s: %v\n", alerta.YouTubeChannel, err)
			continue
		}

		// 5. Decodificar o JSON de resposta
		var ytResp YouTubeSearchResponse
		err = json.NewDecoder(resp.Body).Decode(&ytResp)
		resp.Body.Close() // Fecha o body logo após a leitura

		if err != nil {
			log.Printf("Erro ao decodificar JSON do canal %s: %v\n", alerta.YouTubeChannel, err)
			continue
		}

		// 6. Lógica de Verificação: Se vieram "Items", significa que tem live acontecendo!
		isLiveNow := len(ytResp.Items) > 0

		// 7. Envio de mensagens e atualização do banco
		if isLiveNow && !alerta.IsLive {
			// Entrou em live agora! Envia a mensagem.
			mensagem := fmt.Sprintf("<@&%s> 🔴 **LIVE ON!**\nCola no youtube e vem me assistir: https://www.youtube.com/channel/%s/live", alerta.RoleToMention, alerta.YouTubeChannel)

			_, err := s.ChannelMessageSend(alerta.DiscordChannel, mensagem)
			if err != nil {
				log.Printf("Erro ao enviar mensagem no Discord: %v", err)
			}

			// Atualiza o banco para true para evitar flood no chat [cite: 36]
			err = database.UpdateLiveStatus(context.Background(), alerta.ID, true)
			if err != nil {
				log.Printf("Erro ao atualizar status de live no BD: %v\n", err)
			}

		} else if !isLiveNow && alerta.IsLive {
			// A live acabou. Reseta o status no banco para ficar pronto para a próxima[cite: 39].
			err = database.UpdateLiveStatus(context.Background(), alerta.ID, false)
			if err != nil {
				log.Printf("Erro ao atualizar status de live no BD: %v\n", err)
			}
		}
	}
}
