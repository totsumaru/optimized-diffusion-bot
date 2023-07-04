package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	internal "github.com/techstart35/optimized-diffusion-bot/internal/cmd"
	"github.com/techstart35/optimized-diffusion-bot/internal/errors"
	"github.com/techstart35/optimized-diffusion-bot/internal/id"
)

// スラッシュコマンドを登録します
func RegisterSlashCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != id.ChannelID().TEST {
		return nil
	}

	if err := registerCommand(s, m.GuildID); err != nil {
		return errors.NewError("コマンドを登録できません", err)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "Slashコマンドを追加しました")
	if err != nil {
		return errors.NewError("完了メッセージを送信できません", err)
	}

	return nil
}

// コマンドを登録します
func registerCommand(session *discordgo.Session, guildID string) error {
	commands := []discordgo.ApplicationCommand{
		{
			Name:        internal.SlashCommand_Generate,
			Description: "日本語のプロンプトで画像を生成します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        internal.SlashCommandOptionName_Prompt,
					Description: "出力したい画像を説明してください",
					Required:    true,
					MaxLength:   500,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        internal.SlashCommandOptionName_Seed,
					Description: "シード値（任意）",
					Required:    false,
				},
			},
		},
	}

	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(id.UserID().THIS_BOT, guildID, &command)
		if err != nil {
			return errors.NewError(fmt.Sprintf("コマンドを登録できません Name: %s", command.Name), err)
		}
	}

	return nil
}
