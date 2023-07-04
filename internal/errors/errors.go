// カスタムエラーを提供します
package errors

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/optimized-diffusion-bot/internal/color"
	"github.com/techstart35/optimized-diffusion-bot/internal/id"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

// カスタムエラーです
type Error struct {
	Err      error
	Previous error
}

// エラーメッセージを送信します
func SendErrMsg(s *discordgo.Session, e error, user *discordgo.User) {
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("エラーが発生しました"),
		Description: e.Error(),
		Color:       color.Red,
		Timestamp:   time.Now().Format("2006-01-02T15:04:05+09:00"),
	}

	if user != nil {
		author := &discordgo.MessageEmbedAuthor{
			IconURL: user.AvatarURL(""),
			Name:    user.Username,
		}
		embed.Author = author
	}

	data := &discordgo.MessageSend{
		Content: fmt.Sprintf("<@%s> Target: <@%s>", id.UserID().TOTSUMARU, user.ID),
		Embed:   embed,
	}

	_, err := s.ChannelMessageSendComplex(id.ChannelID().ERR_LOG, data)
	if err != nil {
		log.Println(err)
	}
}

// エラーを新規作成します
//
// 一つ前のエラーを保持しているので、このエラーを再帰的に仕様することで簡易的なスタックトレースを出力できます。
//
// 引数のパターンは以下のとおりです。
//
// 1. エラーメッセージを指定する
//
// NewError("message")
//
// 2. エラーメッセージと一つ前のエラーを指定する
//
// NewError("message", err)
func NewError(msg string, arg ...interface{}) *Error {
	var (
		prev error
	)

	switch len(arg) {
	case 0:
		prev = nil
	case 1:
		switch a := arg[0].(type) {
		case *Error:
			prev = a
		case error:
			prev = a
		}
	}

	_, file, line, _ := runtime.Caller(1)

	if prev != nil {
		return &Error{
			Err:      fmt.Errorf(msg+" file: %s line: %d prev: [%w]", filepath.ToSlash(file), line, prev),
			Previous: prev,
		}
	}

	return &Error{
		Err:      fmt.Errorf(msg+" file: %s line: %d", filepath.ToSlash(file), line),
		Previous: prev,
	}
}

// エラーメッセージを返します
func (e *Error) Error() string {
	return e.Err.Error()
}

// １つ前のエラーを返します
func (e *Error) Unwrap() error {
	return e.Previous
}
