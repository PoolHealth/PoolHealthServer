package common

import "time"

type ActionType uint

const (
	// ActionType enum
	ActionTypeUnknown ActionType = iota
	ActionNet
	ActionBrush
	ActionVacuum
	ActionBackwash
	ActionScumLine
	ActionPumpBasketClean
	ActionSkimmerBasketClean
)

type Action struct {
	Types     []ActionType
	CreatedAt time.Time
}

var ActionTypeNames = map[ActionType]string{
	ActionTypeUnknown:        "Unknown",
	ActionNet:                "Net",
	ActionBrush:              "Brush",
	ActionVacuum:             "Vacuum",
	ActionBackwash:           "Backwash",
	ActionScumLine:           "Scum Line",
	ActionPumpBasketClean:    "Pump Basket Clean",
	ActionSkimmerBasketClean: "Skimmer Basket Clean",
}

var ActionTypeNamesToActionType = map[string]ActionType{
	"Unknown":              ActionTypeUnknown,
	"Net":                  ActionNet,
	"Brush":                ActionBrush,
	"Vacuum":               ActionVacuum,
	"Backwash":             ActionBackwash,
	"Scum Line":            ActionScumLine,
	"Pump Basket Clean":    ActionPumpBasketClean,
	"Skimmer Basket Clean": ActionSkimmerBasketClean,
}
