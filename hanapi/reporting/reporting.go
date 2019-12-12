package reporting

import (
	"fmt"
	"github.com/oliveroneill/slack"
)

// Logger used to log errors or messages
type Logger interface {
	Log(message string)
}

// SlackLogger is an implementation used for logging to slack
type SlackLogger struct {
	Logger
	apiToken string
}

// NewSlackLogger creates a Logger that sends messages to Slack
func NewSlackLogger(apiToken string) *SlackLogger {
	l := new(SlackLogger)
	l.apiToken = apiToken
	return l
}

// Log will log the specified message to Slack under the `hanserver` channel
func (l *SlackLogger) Log(message string) {
	if len(l.apiToken) == 0 {
		fmt.Println("No Slack API token set. Logging disabled")
		return
	}
	// notify through Slack bot
	channelName := "hanserver"
	api := slack.New(l.apiToken)
	params := slack.PostMessageParameters{}
	_, _, err := api.PostMessage(channelName, message, params)
	if err != nil {
		fmt.Println(err)
		return
	}
}
