package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/optimized-diffusion-bot/handler/interaction/generate_jp"
	"github.com/techstart35/optimized-diffusion-bot/internal/cmd"
	"github.com/techstart35/optimized-diffusion-bot/internal/errors"
)

// コマンドが実行された時のハンドラーです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Interaction.Type {
	// メッセージコンポーネント（ボタン）イベント
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		}
	case discordgo.InteractionApplicationCommand:
		name := i.Data.(discordgo.ApplicationCommandInteractionData).Name
		switch name {
		case cmd.SlashCommand_Generate:
			if err := generate_jp.GenerateFromJPPrompt(s, i); err != nil {
				errors.SendErrMsg(s, err, i.Member.User)
			}
			return
		}
	}
}
