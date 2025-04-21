package components

import (
	"context"
	"discord-bot-go/internal/utils"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Estrutura para armazenar a função de cancelamento e o tempo de criação
type cancelEntry struct {
	Cancel    context.CancelFunc
	CreatedAt time.Time
}

// Mapa para armazenar entradas de cancelamento por channelID
var (
	cancelFuncs     = make(map[string]cancelEntry)
	cancelFuncsLock sync.Mutex
)

func init() {
	// Registrar o componente cancel_close
	RegisterComponent(Component{
		CustomID: "cancel_close",
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channelID := strings.TrimPrefix(i.MessageComponentData().CustomID, "cancel_close_")

			// Acionar o cancelamento
			TriggerCancel(channelID)

			// Responder ao usuário confirmando o cancelamento
			utils.RespondEphemeralEmbed(s, i, utils.ResponseOptions{
				Title:       "Fechamento Cancelado",
				Description: "O canal será mantido aberto.",
				Color:       0xFFFF00, // Amarelo para cancelamento
			})

			// Deletar a mensagem de fechamento
			err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				s.ChannelMessageSend(channelID, "Erro ao deletar mensagem de fechamento. O canal permanecerá aberto.")
			}
		},
	})

	// Goroutine para limpeza periódica do mapa cancelFuncs
	go func() {
		// Cria um ticker que dispara a cada 1 minuto
		ticker := time.NewTicker(1 * time.Minute)
		// Garante que o ticker seja parado quando a goroutine terminar
		defer ticker.Stop()

		// Loop que executa a limpeza a cada tick
		for range ticker.C {
			// Envolve a limpeza em uma função com recuperação de panic
			func() {
				// Captura qualquer panic para evitar que a goroutine morra
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in cleanup goroutine: %v", r)
					}
				}()

				// Log do início da limpeza
				log.Printf("Running periodic cleanup of cancelFuncs, current size: %d", len(cancelFuncs))

				// Bloqueia o mapa para acesso seguro
				cancelFuncsLock.Lock()
				defer cancelFuncsLock.Unlock()

				// Itera sobre as entradas do mapa
				for channelID, entry := range cancelFuncs {
					// Só limpa entradas com mais de 30 segundos
					if time.Since(entry.CreatedAt) > 30*time.Second {
						// Cancela o contexto pendente
						entry.Cancel()
						// Remove a entrada do mapa
						delete(cancelFuncs, channelID)
						log.Printf("Cleaned up obsolete cancel func for channel %s", channelID)
					}
				}
			}()
		}
	}()
}

// RegisterCancelFunc registra uma função de cancelamento para um channelID
func RegisterCancelFunc(channelID string, cancel context.CancelFunc) {
	cancelFuncsLock.Lock()
	defer cancelFuncsLock.Unlock()
	cancelFuncs[channelID] = cancelEntry{
		Cancel:    cancel,
		CreatedAt: time.Now(),
	}
	log.Printf("Registered cancel func for channel %s, map size: %d", channelID, len(cancelFuncs))
}

// TriggerCancel aciona o cancelamento para um channelID
func TriggerCancel(channelID string) {
	cancelFuncsLock.Lock()
	defer cancelFuncsLock.Unlock()
	if entry, exists := cancelFuncs[channelID]; exists {
		entry.Cancel()
		delete(cancelFuncs, channelID)
		log.Printf("Triggered cancel for channel %s, map size: %d", channelID, len(cancelFuncs))
	}
}

// ClearCancelFunc remove a função de cancelamento para um channelID
func ClearCancelFunc(channelID string) {
	cancelFuncsLock.Lock()
	defer cancelFuncsLock.Unlock()
	delete(cancelFuncs, channelID)
	log.Printf("Cleared cancel func for channel %s, map size: %d", channelID, len(cancelFuncs))
}
