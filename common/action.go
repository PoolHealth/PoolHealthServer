package common

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
