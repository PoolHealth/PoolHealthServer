// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/PoolHealth/PoolHealthServer/pkg/api/v1/common"
	"github.com/google/uuid"
)

type ChemicalValue interface {
	IsChemicalValue()
}

type AcidChemicalValue struct {
	Type  AcidChemical `json:"type"`
	Value float64      `json:"value"`
}

func (AcidChemicalValue) IsChemicalValue() {}

type AcidChemicalValueInput struct {
	Type  AcidChemical `json:"type"`
	Value float64      `json:"value"`
}

type Action struct {
	Types     []ActionType `json:"types"`
	CreatedAt time.Time    `json:"createdAt"`
}

type AlkalinityChemicalValue struct {
	Type  AlkalinityChemical `json:"type"`
	Value float64            `json:"value"`
}

func (AlkalinityChemicalValue) IsChemicalValue() {}

type AlkalinityChemicalValueInput struct {
	Type  AlkalinityChemical `json:"type"`
	Value float64            `json:"value"`
}

type ChemicalInput struct {
	PoolID     common.ID                       `json:"poolID"`
	Chlorine   []*ChlorineChemicalValueInput   `json:"chlorine,omitempty"`
	Acid       []*AcidChemicalValueInput       `json:"acid,omitempty"`
	Alkalinity []*AlkalinityChemicalValueInput `json:"alkalinity,omitempty"`
}

type Chemicals struct {
	Value     []ChemicalValue `json:"value"`
	CreatedAt time.Time       `json:"createdAt"`
}

type ChlorineChemicalValue struct {
	Type  ChlorineChemical `json:"type"`
	Value float64          `json:"value"`
}

func (ChlorineChemicalValue) IsChemicalValue() {}

type ChlorineChemicalValueInput struct {
	Type  ChlorineChemical `json:"type"`
	Value float64          `json:"value"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CoordinatesInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Measurement struct {
	Chlorine   *float64 `json:"chlorine,omitempty"`
	Ph         *float64 `json:"ph,omitempty"`
	Alkalinity *float64 `json:"alkalinity,omitempty"`
}

type MeasurementRecord struct {
	Measurement *Measurement `json:"measurement"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type Migration struct {
	ID     common.ID       `json:"id"`
	Status MigrationStatus `json:"status"`
}

type Mutation struct {
}

type Pool struct {
	ID       common.ID     `json:"id"`
	Name     string        `json:"name"`
	Volume   float64       `json:"volume"`
	Settings *PoolSettings `json:"settings,omitempty"`
}

type PoolSettings struct {
	Type         PoolType     `json:"type"`
	UsageType    UsageType    `json:"usageType"`
	LocationType LocationType `json:"locationType"`
	Shape        PoolShape    `json:"shape"`
	Coordinates  *Coordinates `json:"coordinates"`
}

type PoolSettingsInput struct {
	Type         PoolType          `json:"type"`
	UsageType    UsageType         `json:"usageType"`
	LocationType LocationType      `json:"locationType"`
	PoolShape    PoolShape         `json:"poolShape"`
	Coordinates  *CoordinatesInput `json:"coordinates"`
}

type Query struct {
}

type Session struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type Subscription struct {
}

type User struct {
	TokenExpiredAt time.Time `json:"tokenExpiredAt"`
	Pools          []*Pool   `json:"pools"`
	id             uuid.UUID `json:"-"`
}

type AcidChemical string

const (
	AcidChemicalHydrochloricAcid AcidChemical = "HydrochloricAcid"
	AcidChemicalSodiumBisulphate AcidChemical = "SodiumBisulphate"
)

var AllAcidChemical = []AcidChemical{
	AcidChemicalHydrochloricAcid,
	AcidChemicalSodiumBisulphate,
}

func (e AcidChemical) IsValid() bool {
	switch e {
	case AcidChemicalHydrochloricAcid, AcidChemicalSodiumBisulphate:
		return true
	}
	return false
}

func (e AcidChemical) String() string {
	return string(e)
}

func (e *AcidChemical) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AcidChemical(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AcidChemical", str)
	}
	return nil
}

func (e AcidChemical) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ActionType string

const (
	ActionTypeNet                ActionType = "Net"
	ActionTypeBrush              ActionType = "Brush"
	ActionTypeVacuum             ActionType = "Vacuum"
	ActionTypeBackwash           ActionType = "Backwash"
	ActionTypeScumLine           ActionType = "ScumLine"
	ActionTypePumpBasketClean    ActionType = "PumpBasketClean"
	ActionTypeSkimmerBasketClean ActionType = "SkimmerBasketClean"
)

var AllActionType = []ActionType{
	ActionTypeNet,
	ActionTypeBrush,
	ActionTypeVacuum,
	ActionTypeBackwash,
	ActionTypeScumLine,
	ActionTypePumpBasketClean,
	ActionTypeSkimmerBasketClean,
}

func (e ActionType) IsValid() bool {
	switch e {
	case ActionTypeNet, ActionTypeBrush, ActionTypeVacuum, ActionTypeBackwash, ActionTypeScumLine, ActionTypePumpBasketClean, ActionTypeSkimmerBasketClean:
		return true
	}
	return false
}

func (e ActionType) String() string {
	return string(e)
}

func (e *ActionType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionType", str)
	}
	return nil
}

func (e ActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type AlkalinityChemical string

const (
	AlkalinityChemicalSodiumBicarbonate AlkalinityChemical = "SodiumBicarbonate"
)

var AllAlkalinityChemical = []AlkalinityChemical{
	AlkalinityChemicalSodiumBicarbonate,
}

func (e AlkalinityChemical) IsValid() bool {
	switch e {
	case AlkalinityChemicalSodiumBicarbonate:
		return true
	}
	return false
}

func (e AlkalinityChemical) String() string {
	return string(e)
}

func (e *AlkalinityChemical) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AlkalinityChemical(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AlkalinityChemical", str)
	}
	return nil
}

func (e AlkalinityChemical) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ChlorineChemical string

const (
	ChlorineChemicalCalciumHypochlorite65Percent ChlorineChemical = "CalciumHypochlorite65Percent"
	ChlorineChemicalSodiumHypochlorite12Percent  ChlorineChemical = "SodiumHypochlorite12Percent"
	ChlorineChemicalSodiumHypochlorite14Percent  ChlorineChemical = "SodiumHypochlorite14Percent"
	ChlorineChemicalTCCA90PercentTablets         ChlorineChemical = "TCCA90PercentTablets"
	ChlorineChemicalMultiActionTablets           ChlorineChemical = "MultiActionTablets"
	ChlorineChemicalTCCA90PercentGranules        ChlorineChemical = "TCCA90PercentGranules"
	ChlorineChemicalDichlor65Percent             ChlorineChemical = "Dichlor65Percent"
)

var AllChlorineChemical = []ChlorineChemical{
	ChlorineChemicalCalciumHypochlorite65Percent,
	ChlorineChemicalSodiumHypochlorite12Percent,
	ChlorineChemicalSodiumHypochlorite14Percent,
	ChlorineChemicalTCCA90PercentTablets,
	ChlorineChemicalMultiActionTablets,
	ChlorineChemicalTCCA90PercentGranules,
	ChlorineChemicalDichlor65Percent,
}

func (e ChlorineChemical) IsValid() bool {
	switch e {
	case ChlorineChemicalCalciumHypochlorite65Percent, ChlorineChemicalSodiumHypochlorite12Percent, ChlorineChemicalSodiumHypochlorite14Percent, ChlorineChemicalTCCA90PercentTablets, ChlorineChemicalMultiActionTablets, ChlorineChemicalTCCA90PercentGranules, ChlorineChemicalDichlor65Percent:
		return true
	}
	return false
}

func (e ChlorineChemical) String() string {
	return string(e)
}

func (e *ChlorineChemical) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChlorineChemical(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChlorineChemical", str)
	}
	return nil
}

func (e ChlorineChemical) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type LocationType string

const (
	LocationTypeIndoor  LocationType = "Indoor"
	LocationTypeOutdoor LocationType = "Outdoor"
)

var AllLocationType = []LocationType{
	LocationTypeIndoor,
	LocationTypeOutdoor,
}

func (e LocationType) IsValid() bool {
	switch e {
	case LocationTypeIndoor, LocationTypeOutdoor:
		return true
	}
	return false
}

func (e LocationType) String() string {
	return string(e)
}

func (e *LocationType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = LocationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid LocationType", str)
	}
	return nil
}

func (e LocationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MigrationStatus string

const (
	MigrationStatusUnknown MigrationStatus = "Unknown"
	MigrationStatusPending MigrationStatus = "Pending"
	MigrationStatusDone    MigrationStatus = "Done"
	MigrationStatusFailed  MigrationStatus = "Failed"
)

var AllMigrationStatus = []MigrationStatus{
	MigrationStatusUnknown,
	MigrationStatusPending,
	MigrationStatusDone,
	MigrationStatusFailed,
}

func (e MigrationStatus) IsValid() bool {
	switch e {
	case MigrationStatusUnknown, MigrationStatusPending, MigrationStatusDone, MigrationStatusFailed:
		return true
	}
	return false
}

func (e MigrationStatus) String() string {
	return string(e)
}

func (e *MigrationStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MigrationStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MigrationStatus", str)
	}
	return nil
}

func (e MigrationStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

var AllOrder = []Order{
	OrderAsc,
	OrderDesc,
}

func (e Order) IsValid() bool {
	switch e {
	case OrderAsc, OrderDesc:
		return true
	}
	return false
}

func (e Order) String() string {
	return string(e)
}

func (e *Order) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Order(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Order", str)
	}
	return nil
}

func (e Order) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PoolShape string

const (
	PoolShapeRectangle PoolShape = "Rectangle"
	PoolShapeCircle    PoolShape = "Circle"
	PoolShapeOval      PoolShape = "Oval"
	PoolShapeKidney    PoolShape = "Kidney"
	PoolShapeL         PoolShape = "L"
	PoolShapeT         PoolShape = "T"
	PoolShapeFreeForm  PoolShape = "FreeForm"
)

var AllPoolShape = []PoolShape{
	PoolShapeRectangle,
	PoolShapeCircle,
	PoolShapeOval,
	PoolShapeKidney,
	PoolShapeL,
	PoolShapeT,
	PoolShapeFreeForm,
}

func (e PoolShape) IsValid() bool {
	switch e {
	case PoolShapeRectangle, PoolShapeCircle, PoolShapeOval, PoolShapeKidney, PoolShapeL, PoolShapeT, PoolShapeFreeForm:
		return true
	}
	return false
}

func (e PoolShape) String() string {
	return string(e)
}

func (e *PoolShape) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PoolShape(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PoolShape", str)
	}
	return nil
}

func (e PoolShape) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PoolType string

const (
	PoolTypeInfinity PoolType = "Infinity"
	PoolTypeOverflow PoolType = "Overflow"
	PoolTypeSkimmer  PoolType = "Skimmer"
)

var AllPoolType = []PoolType{
	PoolTypeInfinity,
	PoolTypeOverflow,
	PoolTypeSkimmer,
}

func (e PoolType) IsValid() bool {
	switch e {
	case PoolTypeInfinity, PoolTypeOverflow, PoolTypeSkimmer:
		return true
	}
	return false
}

func (e PoolType) String() string {
	return string(e)
}

func (e *PoolType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PoolType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PoolType", str)
	}
	return nil
}

func (e PoolType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UsageType string

const (
	UsageTypeCommunity UsageType = "Community"
	UsageTypePrivate   UsageType = "Private"
	UsageTypeHoliday   UsageType = "Holiday"
)

var AllUsageType = []UsageType{
	UsageTypeCommunity,
	UsageTypePrivate,
	UsageTypeHoliday,
}

func (e UsageType) IsValid() bool {
	switch e {
	case UsageTypeCommunity, UsageTypePrivate, UsageTypeHoliday:
		return true
	}
	return false
}

func (e UsageType) String() string {
	return string(e)
}

func (e *UsageType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UsageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UsageType", str)
	}
	return nil
}

func (e UsageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
