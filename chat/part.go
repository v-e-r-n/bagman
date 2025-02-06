package chat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/v-e-r-n/bagman"
	"github.com/v-e-r-n/bagman/utils"
)

var _ bagman.Part = (*Part)(nil)

type Part struct {
	data map[string]any
}

func (p *Part) Kind() string {
	kind, _ := utils.Get[string](p.data, "type")

	return strings.ToLower(kind)
}

func (p *Part) Content() (any, error) {
	return p.Export()
}

func (p *Part) Update(content any) error {
	maybe, ok := content.(map[string]any)
	if !ok {
		return fmt.Errorf("content is not an object")
	}

	p.data = maybe

	return nil
}

func (p *Part) Import(blob any) error {
	if data, err := utils.Objectify(blob); err == nil {
		p.data = data

		return nil
	}

	data := map[string]any{
		"type": "text",
		"text": blob,
	}

	p.data = data

	return nil
}

func (p *Part) Export() (any, error) {
	return p.data, nil
}

func (p *Part) MarshalJSON() ([]byte, error) {
	data, err := p.Export()
	if err != nil {
		return nil, err
	}

	return json.Marshal(data)
}

func (p *Part) UnmarshalJSON(data []byte) error {
	raw := make(map[string]any)

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	return p.Import(raw)
}

func (p *Part) String() string {
	blah, _ := p.MarshalJSON()

	return string(blah)
}
