package proxy

import (
	"encoding/json"
	"fmt"

	"github.com/v-e-r-n/bagman"
)

var _ bagman.Request = (*GenericRequest)(nil)

func NewGenericRequest() *GenericRequest {
	return &GenericRequest{
		data: make(map[string]any),
	}
}

type GenericRequest struct {
	data map[string]any
}

func (r *GenericRequest) Import(any) error {
	return fmt.Errorf("not implemented")
}

func (r *GenericRequest) Export() (any, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *GenericRequest) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *GenericRequest) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("not implemented")
}

func (r *GenericRequest) Kind() string {
	return "GenericRequest"
}

func (r *GenericRequest) MessagesFor(string) []bagman.Message {
	return nil
}

func requestify(upstreamRequest bagman.Request, mod RequestModifier, logger bagman.Logger) BodyMod {
	return func(body []byte) ([]byte, error) {
		// upstreamRequest := NewGenericRequest()
		err := json.Unmarshal(body, &upstreamRequest)
		if err != nil {
			logger.Error("failed to unmarshal request", "error", err)
			return nil, err
		}

		mod(upstreamRequest)

		proxyBody, err := json.Marshal(upstreamRequest)
		if err != nil {
			logger.Error("failed to marshal upstream request", "error", err)
			return nil, err
		}

		return proxyBody, nil
	}
}
