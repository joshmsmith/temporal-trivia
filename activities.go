package triviagame

import (
	"context"

	"github.com/ktenzer/triviagame/resources"
	"github.com/sashabaranov/go-openai"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

func TriviaQuestionActivity(ctx context.Context, input resources.ActivityInput) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("TriviaQuestionActivity")

	client := openai.NewClient(input.Key)
	messages := make([]openai.ChatCompletionMessage, 0)

	result, err := resources.SendChatGptRequest(client, messages, input.Question)

	if err != nil {
		return result, temporal.NewApplicationError("ChatGPT", "request", err)
	}

	return result, nil
}

func ScoreTotalActivity(ctx context.Context, scoreMap map[string]int) (resources.ActivityScoreOutput, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("ScoreTotalActivity")

	var highScore int
	var winners []string
	for user, score := range scoreMap {
		if score >= highScore {
			highScore = score
			winners = append(winners, user)
		}
	}

	output := &resources.ActivityScoreOutput{
		Winners:   winners,
		HighScore: highScore,
	}

	return *output, nil
}
