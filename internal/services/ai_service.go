package services

import (
	"ac-ai/internal/config"
	"ac-ai/internal/models"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		client: openai.NewClient(cfg.OpenAIAPIKey),
	}
}

// 1. Извлечение суммы (без изменений)
func (s *AIService) ParseAmountFromQuery(ctx context.Context, query string) (float64, error) {
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Ты - парсер. Извлеки число (сумму) из запроса. Ответь ТОЛЬКО числом (например '15000000'). Если числа нет, ответь '0'.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
			Temperature: 0,
		},
	)

	if err != nil {
		return 0, err
	}

	amountStr := resp.Choices[0].Message.Content
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	// ** Исправляем парсинг для больших чисел **
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// 2. "Теплый" анализ (ИСПРАВЛЕННЫЙ ПРОМПТ)
func (s *AIService) GetAIAnalysis(ctx context.Context, scoreData *ColdScoreResult, profile *models.FinancialProfile) (string, error) {

	scoreDataBytes, _ := json.MarshalIndent(scoreData, "", "  ")

	systemPrompt := `
	Ты - AI-ассистент банка, вежливый и профессиональный кредитный аналитик.
	Твоя задача - проанализировать РЕЗУЛЬТАТЫ СКОРИНГА и дать клиенту развернутый,
	человекопонятный ответ.
	
	НЕ ПОКАЗЫВАЙ клиенту баллы, DTI или 'breakdown'.
	Твой ответ должен основываться на поле 'decision', 'recommendations' и 'recommendedMaxAmount'.
	
	Вот РЕЗУЛЬТАТЫ СКОРИНГА (это главный документ):
	` + string(scoreDataBytes) + `
	
	Твоя задача - выполнить следующий алгоритм:

	1.  **Посмотри на 'decision'.**

	2.  **Если 'decision' == "APPROVED":**
		* Поздравь клиента.
		* Сообщи, что заявка на ` + strconv.FormatFloat(scoreData.RequestedAmount, 'f', 0, 64) + ` тг предварительно одобрена.

	3.  **Если 'decision' == "MANUAL_REVIEW":**
		* Сообщи, что заявка на ` + strconv.FormatFloat(scoreData.RequestedAmount, 'f', 0, 64) + ` тг отправлена на ручное рассмотрение.
		* **Объясни причину:** Посмотри на 'recommendations'. Вежливо перечисли 1-2 основные причины (например, "из-за недавних просрочек" или "из-за высокого стажа").
		* **Проверь сумму:** Если 'requestedAmount' > 'recommendedMaxAmount', обязательно добавь: "В частности, запрошенная вами сумма может быть слишком высокой для вашего текущего дохода. Возможно, наш менеджер предложит вам скорректированную сумму."

	4.  **Если 'decision' == "DENIED":**
		* Вежливо сообщи об отказе по заявке на ` + strconv.FormatFloat(scoreData.RequestedAmount, 'f', 0, 64) + ` тг.
		* **ОБЯЗАТЕЛЬНО объясни главную причину:**
			* **Сценарий 1: Сумма слишком велика (ЭТО ГЛАВНЫЙ СЦЕНАРИЙ ДЛЯ 50 МЛРД).**
				* Проверь, если 'requestedAmount' > 'recommendedMaxAmount'.
				* Если это так, скажи: "К сожалению, в кредите отказано. Основная причина - запрошенная сумма ( ` + strconv.FormatFloat(scoreData.RequestedAmount, 'f', 0, 64) + ` тг) слишком велика для вашего текущего уровня подтвержденного дохода."
				* **!!ДАЙ АЛЬТЕРНАТИВУ!!:** "На основе вашего профиля, мы могли бы быстро одобрить вам сумму в размере **` + strconv.FormatFloat(scoreData.RecommendedMaxAmount, 'f', 0, 64) + ` тг**. Вы можете подать повторную заявку на эту сумму."
			* **Сценарий 2: Плохая кредитная история или другие факторы (сумма в порядке).**
				* Если 'requestedAmount' <= 'recommendedMaxAmount' (т.е. дело не в сумме), посмотри на 'recommendations'.
				* Скажи: "К сожалению, в кредите отказано. Основные причины: " (и перечисли 1-2 пункта из 'recommendations', например, "наличие серьезных просрочек в кредитной истории" или "высокая текущая долговая нагрузка").
				* **ДАЙ СОВЕТ:** "Мы рекомендуем вам [совет на основе причины, например: 'улучшить вашу кредитную историю, закрыв текущие просрочки'] и попробовать подать заявку через несколько месяцев."

	Используй вежливый и заботливый тон.
	`

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Сформулируй ответ для клиента.",
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
