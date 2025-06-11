package graphql

import (
	"errors"

	"github.com/PoolHealth/PoolHealthServer/common"
)

var ErrUnknownActionType = errors.New("unknown action type")

func (a ActionType) ToCommon() common.ActionType {
	switch a {
	case ActionTypeNet:
		return common.ActionNet
	case ActionTypeBrush:
		return common.ActionBrush
	case ActionTypeVacuum:
		return common.ActionVacuum
	case ActionTypeBackwash:
		return common.ActionBackwash
	case ActionTypeScumLine:
		return common.ActionScumLine
	case ActionTypePumpBasketClean:
		return common.ActionPumpBasketClean
	case ActionTypeSkimmerBasketClean:
		return common.ActionSkimmerBasketClean
	default:
		return common.ActionTypeUnknown
	}
}

func ActionTypeFromCommon(a common.ActionType) (ActionType, error) {
	switch a {
	case common.ActionNet:
		return ActionTypeNet, nil
	case common.ActionBrush:
		return ActionTypeBrush, nil
	case common.ActionVacuum:
		return ActionTypeVacuum, nil
	case common.ActionBackwash:
		return ActionTypeBackwash, nil
	case common.ActionScumLine:
		return ActionTypeScumLine, nil
	case common.ActionPumpBasketClean:
		return ActionTypePumpBasketClean, nil
	case common.ActionSkimmerBasketClean:
		return ActionTypeSkimmerBasketClean, nil
	default:
		return "", ErrUnknownActionType
	}
}

func ActionFromCommon(a common.Action) (*Action, error) {
	action := &Action{
		Types:     make([]ActionType, len(a.Types)),
		CreatedAt: a.CreatedAt,
	}

	for i, k := range a.Types {
		actionType, err := ActionTypeFromCommon(k)
		if err != nil {
			return nil, err
		}

		action.Types[i] = actionType
	}

	return action, nil
}
