package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/optimized-diffusion-bot/handler/interaction"
	"github.com/techstart35/optimized-diffusion-bot/handler/message"
)

// メッセージが作成された時のハンドラです
func Handler(s *discordgo.Session) {
	s.AddHandler(message.MessageCreateHandler)
	s.AddHandler(interaction.InteractionCreateHandler)
}
