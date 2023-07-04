package interaction

import (
	"github.com/bwmarrin/discordgo"
)

// コマンドが実行された時のハンドラーです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Interaction.Type {
	// メッセージコンポーネント（ボタン）イベント
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		}
	}
}
