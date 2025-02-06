package chat

import (
	"encoding/json"
	"fmt"

	"github.com/v-e-r-n/bagman"
	"github.com/v-e-r-n/bagman/utils"
)

var _ bagman.Message = (*Message)(nil)

type Message struct {
	Role  string
	Parts []bagman.Part
	data  map[string]any
}

func (m *Message) Kind() string {
	return m.Role
}

func (m *Message) Content() (any, error) {
	return m.Parts, nil
}

func (m *Message) AllParts() []bagman.Part {
	return m.Parts
}

func (m *Message) PartsOf(kind string) []bagman.Part {
	parts := make([]bagman.Part, 0)

	for _, part := range m.Parts {
		if part.Kind() == kind {
			parts = append(parts, part)
		}
	}

	return parts
}

func extractSimpleContent(blob any) ([]bagman.Part, error) {
	content, ok := blob.(string)
	if !ok {
		return nil, fmt.Errorf("content is not a string")
	}

	parts := []bagman.Part{
		&Part{
			data: map[string]any{
				"type": "text",
				"text": content,
			},
		},
	}

	return parts, nil
}

func extractComplexContent(blob any) ([]bagman.Part, error) {
	content, ok := blob.([]any)
	if !ok {
		return nil, fmt.Errorf("content is not an array")
	}

	parts := make([]bagman.Part, 0)
	for _, part := range content {
		p := &Part{}
		err := p.Import(part)
		if err != nil {
			return nil, err
		}

		parts = append(parts, p)
	}

	return parts, nil
}

func extractContent(blob any) ([]bagman.Part, error) {
	if parts, err := extractSimpleContent(blob); err == nil {
		return parts, nil
	}

	return extractComplexContent(blob)
}

func (m *Message) Import(blob any) error {
	data, err := utils.Objectify(blob)
	if err != nil {
		fmt.Println("well shit 1")
		return err
	}

	role, err := utils.Get[string](data, "role")
	if err != nil {
		fmt.Println("well shit 2")
		return err
	}

	content, found := data["content"]
	if !found {
		return fmt.Errorf("message has no content")
	}

	parts, err := extractContent(content)
	if err != nil {
		return err
	}

	m.Role = role
	m.Parts = parts
	m.data = data

	return nil
}

func (m *Message) Export() (any, error) {
	data := m.data
	data["role"] = m.Role
	data["content"] = m.Parts

	return data, nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	data, err := m.Export()
	if err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (m *Message) UnmarshalJSON(data []byte) error {
	blob := make(map[string]any)

	err := json.Unmarshal(data, &blob)
	if err != nil {
		return err
	}

	return m.Import(blob)
}

func (m *Message) String() string {
	blah, _ := m.MarshalJSON()

	return string(blah)
}
