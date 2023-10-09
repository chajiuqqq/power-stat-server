package service

import (
	"context"
	"mqtt-wx-forward/types"
)

func (s *Service) Report(ctx context.Context, msg string) (string, error) {
	switch msg {
	case types.MsgStat:
		return types.MsgStat, nil
	case types.MsgStatMonthly:
		return types.MsgStatMonthly, nil
	default:
		return "", nil
	}
}
