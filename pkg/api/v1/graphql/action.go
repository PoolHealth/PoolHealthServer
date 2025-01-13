package graphql

import (
	"errors"

	"github.com/PoolHealth/PoolHealthServer/common"
)

func (a Action) ToCommon() common.ActionType {
	switch a {
	case ActionNet:
		return common.ActionNet
	case ActionBrush:
		return common.ActionBrush
	case ActionVacuum:
		return common.ActionVacuum
	case ActionBackwash:
		return common.ActionBackwash
	case ActionScumLine:
		return common.ActionScumLine
	case ActionPumpBasketClean:
		return common.ActionPumpBasketClean
	case ActionSkimmerBasketClean:
		return common.ActionSkimmerBasketClean
	default:
		return common.ActionTypeUnknown
	}
}

func ActionFromCommon(a common.ActionType) (Action, error) {
	switch a {
	case common.ActionNet:
		return ActionNet, nil
	case common.ActionBrush:
		return ActionBrush, nil
	case common.ActionVacuum:
		return ActionVacuum, nil
	case common.ActionBackwash:
		return ActionBackwash, nil
	case common.ActionScumLine:
		return ActionScumLine, nil
	case common.ActionPumpBasketClean:
		return ActionPumpBasketClean, nil
	case common.ActionSkimmerBasketClean:
		return ActionSkimmerBasketClean, nil
	default:
		return "", errors.New("unknown action type")
	}
}
