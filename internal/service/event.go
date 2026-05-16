package service

import (
	"context"
	"encoding/json"
)

func (s *Service) HandleEvent(ctx context.Context, msg []byte) error {
	var data map[string]interface{}

	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	if _, ok := data["event_type"]; ok {
		return s.SendTaskCountMessage(ctx)
	}

	return nil
}
