package chat

import (
	"encoding/json"

	"github.com/v-e-r-n/baagman"
	"github.com/v-e-r-n/baagman/utils"
)

var _ bagman.Object = (*CompletionRequest)(nil)

func NewCompletionRequest() *CompletionRequest {
	return &CompletionRequest{
		data: make(map[string]any),
	}
}

type CompletionRequest struct {
	Messages []*Message     `json:"messages,required"`
	Model    string         `json:"model,required"`
	Stream   *bool          `json:"stream,omitempty"`
	data     map[string]any `json:"-"`
}

func (s *CompletionRequest) String() string {
	data, _ := s.MarshalJSON()

	return string(data)
}

func (s *CompletionRequest) Import(blob any) error {
	data, err := utils.Objectify(blob)
	if err != nil {
		return err
	}

	model, err := utils.Get[string](data, "model")
	if err != nil {
		return err
	}

	messageBlob, err := utils.Get[[]any](data, "messages")
	if err != nil {
		return err
	}

	messages := make([]*Message, 0)
	for _, blob := range messageBlob {
		message := &Message{}
		err = message.Import(blob)
		if err != nil {
			return err
		}

		messages = append(messages, message)
	}

	if stream, found, err := utils.Find[bool](data, "stream"); found {
		if err != nil {
			return err
		}

		s.Stream = &stream
	}

	s.Model = model
	s.Messages = messages
	s.data = data

	return nil

}

func (s *CompletionRequest) Export() (any, error) {
	messages := make([]any, 0)
	for _, message := range s.Messages {
		blob, err := message.Export()
		if err != nil {
			return nil, err
		}

		messages = append(messages, blob)
	}

	s.data["messages"] = messages
	s.data["model"] = s.Model

	return s.data, nil
}

func (c *CompletionRequest) UnmarshalJSON(data []byte) error {
	var parsed map[string]interface{}

	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}

	return c.Import(parsed)
}

func (c *CompletionRequest) MessagesFor(role string) []bagman.Message {
	messages := make([]bagman.Message, 0)

	for _, m := range c.Messages {
		if m.Kind() == role {
			messages = append(messages, m)
		}
	}

	return messages
}

func (c *CompletionRequest) MarshalJSON() ([]byte, error) {
	output, err := c.Export()
	if err != nil {
		return nil, err
	}

	return json.Marshal(output)
}

func (c *CompletionRequest) IsStreaming() bool {
	return c.Stream != nil && *c.Stream
}

func (c *CompletionRequest) Kind() string {
	return "completion_request"
}
