package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/optimized-diffusion-bot/handler/message/command"
	"github.com/techstart35/optimized-diffusion-bot/internal/cmd"
	"github.com/techstart35/optimized-diffusion-bot/internal/errors"
)

// メッセージが作成された時のハンドラーです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case cmd.CMD_AddSlashCommand:
		if err := command.RegisterSlashCommand(s, m); err != nil {
			errors.SendErrMsg(s, err, m.Author)
		}
		return
	}
}
