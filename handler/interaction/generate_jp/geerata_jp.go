package generate_jp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"github.com/techstart35/optimized-diffusion-bot/internal/cmd"
	"github.com/techstart35/optimized-diffusion-bot/internal/errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	PathTxt2Img = "/sdapi/v1/txt2img"
)

// 日本語プロンプトの処理を実行します
func GenerateFromJPPrompt(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var (
		chatGPTPrompt string
		seed          = -1
	)

	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case cmd.SlashCommandOptionName_Prompt:
			chatGPTPrompt = opt.Value.(string)
		case cmd.SlashCommandOptionName_Seed:
			seed = int(opt.Value.(float64))
		}
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "作成しています。1分ほどお待ちください...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	// ChatGPTへプロンプト作成依頼
	imagePrompt, err := TranslateJPToPrompt(chatGPTPrompt)
	if err != nil {
		return errors.NewError("GPTでプロンプトを生成できません", err)
	}

	// 作成
	res, err := generateByStableDiffusion(imagePrompt, seed)
	if err != nil {
		return errors.NewError("stable-diffusionで画像を生成できません", err)
	}

	files := make([]*discordgo.File, 0)
	for _, v := range res {
		files = append(files, v.File)
	}

	contentTmpl := `
画像を生成しました。
seed: 
%d｜%d｜%d
%d｜%d｜%d
`

	params := &discordgo.WebhookParams{
		Content: fmt.Sprintf(
			contentTmpl,
			res[0].Seed,
			res[1].Seed,
			res[2].Seed,
			res[3].Seed,
			res[4].Seed,
			res[5].Seed,
		),
		Flags: discordgo.MessageFlagsEphemeral,
		Files: files,
	}
	if _, err = s.FollowupMessageCreate(i.Interaction, true, params); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}

// 日本語の文字列をGPTでプロンプトに変換します
func TranslateJPToPrompt(jpText string) (string, error) {
	requestTmpl := `
以下の日本語のテキストを、Stable Diffusionで理解できるようなプロンプト(英語)に変換してください。
- 文節カンマ区切りで出力
- 最後はピリオド無し
- 返答は、プロンプトのみ。それ以外の言葉は返信しないでください。

---
%s
`

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(requestTmpl, jpText),
				},
			},
		},
	)

	if err != nil {
		return "", errors.NewError("GPTにリクエストできません", err)
	}

	res := resp.Choices[0].Message.Content

	return res, nil
}

// データを表す構造体
type Data struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt"`
	Steps          int    `json:"steps"`
	BatchSize      int    `json:"batch_size"`
	Seed           int    `json:"seed"`
}

type ResponseData struct {
	Images []string `json:"images"`
	Info   string   `json:"info"`
}

type Info struct {
	AllSeeds []int `json:"all_seeds"`
}

type Res struct {
	File *discordgo.File
	Seed int
}

// stable-diffusionで画像を生成します
func generateByStableDiffusion(prompt string, seed int) ([]Res, error) {
	promptTmpl := "(8k, RAW photo, best quality, masterpiece:1.2), ultra detailed, ultra high res, professional photograph, extremely detailed beautiful girl, extremely detailed face, extremely detailed eyes, extremely detailed skin, extremely detailed fingers, extremely detailed nose, extremely detailed mouth, perfect anatomy,Photo of (Pretty Japanese woman),(%s)"
	// データを作成
	data := &Data{
		Prompt:         fmt.Sprintf(promptTmpl, prompt),
		NegativePrompt: "EasyNegative, (worst quality:2), (low quality:2), (normal quality:2), lowers, normal quality, (monochrome:1.2), (grayscale:1.2),skin spots, skin blemishes, age spot, ugly face, glans, fat, missing fingers, extra fingers, extra arms, extra legs, watermark, text, error, blurry, jpeg artifacts, cropped, bad anatomy, double navel, muscle, nsfw, nude,((selfie))",
		Steps:          20,
		BatchSize:      6,
		Seed:           seed,
	}

	// JSONにエンコード
	j, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	// POSTリクエストを作成
	url := os.Getenv("STABLE_DIFFUSION_URL") + PathTxt2Img
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, errors.NewError("httpリクエストを作成できません", err)
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントを作成しリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewError("httpリクエストを送信できません", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewError("レスポンスを読み込めません", err)
	}

	// JSONを解析するための変数を準備
	var responseData ResponseData

	// JSONを解析
	if err = json.Unmarshal(body, &responseData); err != nil {
		return nil, errors.NewError("JSONを構造体に変換できません", err)
	}

	var info Info

	// infoをもう一度Unmarshal
	if err = json.Unmarshal([]byte(responseData.Info), &info); err != nil {
		return nil, errors.NewError("JSONを構造体に変換できません", err)
	}

	res := make([]Res, 0)

	// 各画像をデコードして保存
	for i, img := range responseData.Images {
		// base64デコード
		b, err := base64.StdEncoding.DecodeString(img)
		if err != nil {
			return nil, errors.NewError("base64にデコードできません", err)
		}

		file := &discordgo.File{
			Name:   "image.png",
			Reader: strings.NewReader(string(b)),
		}

		r := Res{
			File: file,
			Seed: info.AllSeeds[i],
		}

		res = append(res, r)
	}

	return res, nil
}
