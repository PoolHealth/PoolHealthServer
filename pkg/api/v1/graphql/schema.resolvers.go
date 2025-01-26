package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.62

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	null "github.com/guregu/null/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"

	rootCommon "github.com/PoolHealth/PoolHealthServer/common"
	authPkg "github.com/PoolHealth/PoolHealthServer/internal/services/auth"
	"github.com/PoolHealth/PoolHealthServer/pkg/api/v1/common"
)

// AuthApple is the resolver for the authApple field.
func (r *mutationResolver) AuthApple(ctx context.Context, appleCode string, deviceID common.ID) (*Session, error) {
	session, err := r.auth.Auth(ctx, appleCode)
	if err != nil {
		return nil, gqlerror.Wrap(err)
	}
	return &Session{
		Token:     session.JWT,
		ExpiredAt: session.ExpiredAt,
	}, nil
}

// AddPool is the resolver for the addPool field.
func (r *mutationResolver) AddPool(ctx context.Context, name string, volume float64) (*Pool, error) {
	user, err := authPkg.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	pool, err := r.poolData.Create(ctx, user.ID, &rootCommon.PoolData{
		Name:   name,
		Volume: volume,
	})
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return PoolFromCommon(pool), nil
}

// AddMeasurement is the resolver for the addMeasurement field.
func (r *mutationResolver) AddMeasurement(ctx context.Context, poolID common.ID, chlorine *float64, ph *float64, alkalinity *float64) (*MeasurementRecord, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.measurementHistory.CreateMeasurement(ctx, rootCommon.Measurement{
		PoolID:     uuid.UUID(poolID),
		Chlorine:   null.FloatFromPtr(chlorine),
		PH:         null.FloatFromPtr(ph),
		Alkalinity: null.FloatFromPtr(alkalinity),
	})
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return MeasurementRecordFromCommon(res), nil
}

// DeleteMeasurement is the resolver for the deleteMeasurement field.
func (r *mutationResolver) DeleteMeasurement(ctx context.Context, poolID common.ID, createdAt time.Time) (bool, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return false, err
	}

	return r.measurementHistory.DeleteMeasurement(ctx, uuid.UUID(poolID), createdAt)
}

// AddChemicals is the resolver for the addChemicals field.
func (r *mutationResolver) AddChemicals(ctx context.Context, input ChemicalInput) (*Chemicals, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(input.PoolID)); err != nil {
		return nil, err
	}

	res, err := r.additivesHistory.CreateChemicals(ctx, uuid.UUID(input.PoolID), input.ToCommonProducts())
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return ChemicalsFromCommon(res), nil
}

// DeleteChemicals is the resolver for the deleteChemicals field.
func (r *mutationResolver) DeleteChemicals(ctx context.Context, poolID common.ID, createdAt time.Time) (bool, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return false, err
	}

	return r.additivesHistory.DeleteChemicals(ctx, uuid.UUID(poolID), createdAt)
}

// LogActions is the resolver for the logActions field.
func (r *mutationResolver) LogActions(ctx context.Context, poolID common.ID, action []ActionType) (*time.Time, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	act := make([]rootCommon.ActionType, len(action))
	for i, a := range action {
		act[i] = a.ToCommon()
	}

	createdAt, err := r.actions.LogActions(ctx, uuid.UUID(poolID), act)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return &createdAt, nil
}

// DeleteActionsLog is the resolver for the deleteActionsLog field.
func (r *mutationResolver) DeleteActionsLog(ctx context.Context, poolID common.ID, createdAt time.Time) (bool, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return false, err
	}

	return r.actions.DeleteAction(ctx, uuid.UUID(poolID), createdAt)
}

// UpdatePoolSettings is the resolver for the updatePoolSettings field.
func (r *mutationResolver) UpdatePoolSettings(ctx context.Context, poolID common.ID, settings PoolSettingsInput) (*PoolSettings, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.poolSettingsManager.SetSettings(ctx, uuid.UUID(poolID), settings.ToCommon())
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return PoolSettingsFromCommon(res)
}

// MigrateFromSheet is the resolver for the migrateFromSheet field.
func (r *mutationResolver) MigrateFromSheet(ctx context.Context, sheetLink string) (common.ID, error) {
	user, err := authPkg.GetUser(ctx)
	if err != nil {
		return common.ID{}, castGQLError(ctx, err)
	}

	sheetLinkPart := strings.Replace(sheetLink, "https://docs.google.com/spreadsheets/d/", "", 1)
	sheetLinkSliced := strings.Split(sheetLinkPart, "/")
	if len(sheetLinkSliced) < 1 {
		return common.ID{}, gqlerror.Errorf("invalid sheet link")
	}

	id := r.migrator.Migrate(ctx, user.ID, sheetLinkSliced[0])

	return common.ID(id), nil
}

// Settings is the resolver for the settings field.
func (r *poolResolver) Settings(ctx context.Context, obj *Pool) (*PoolSettings, error) {
	settings, err := r.poolSettingsManager.GetSettings(ctx, uuid.UUID(obj.ID))
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return PoolSettingsFromCommon(settings)
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	user, err := authPkg.GetUser(ctx)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}
	return &User{TokenExpiredAt: user.Session.ExpiredAt, id: user.ID}, nil
}

// Pools is the resolver for the pools field.
func (r *queryResolver) Pools(ctx context.Context) ([]*Pool, error) {
	user, err := authPkg.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	pools, err := r.poolData.List(ctx, user.ID)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	result := make([]*Pool, len(pools))
	for i, pool := range pools {
		result[i] = PoolFromCommon(&pool)
	}

	return result, nil
}

// EstimateMeasurement is the resolver for the estimateMeasurement field.
func (r *queryResolver) EstimateMeasurement(ctx context.Context, input ChemicalInput) (*Measurement, error) {
	poolID := uuid.UUID(input.PoolID)
	if err := r.checkAccessToPool(ctx, uuid.UUID(input.PoolID)); err != nil {
		return nil, err
	}

	res, err := r.estimator.EstimateMeasurement(ctx, poolID, input.ToCommonProducts(), r.getMeasurements(ctx))
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return MeasurementFromCommon(res), nil
}

// DemandMeasurement is the resolver for the demandMeasurement field.
func (r *queryResolver) DemandMeasurement(ctx context.Context, poolID common.ID) (*Measurement, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.estimator.DemandMeasurement(ctx, uuid.UUID(poolID))
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return MeasurementFromCommon(res), nil
}

// HistoryOfMeasurement is the resolver for the historyOfMeasurement field.
func (r *queryResolver) HistoryOfMeasurement(ctx context.Context, poolID common.ID, order Order, offset *int, limit *int) ([]*MeasurementRecord, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.measurementHistory.QueryMeasurement(ctx, uuid.UUID(poolID), order.ToCommon(), offset, limit)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	result := make([]*MeasurementRecord, len(res))
	for i, m := range res {
		result[i] = MeasurementRecordFromCommon(m)
	}

	return result, nil
}

// HistoryOfAdditives is the resolver for the historyOfAdditives field.
func (r *queryResolver) HistoryOfAdditives(ctx context.Context, poolID common.ID, order Order, offset *int, limit *int) ([]*Chemicals, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.additivesHistory.QueryAdditives(ctx, uuid.UUID(poolID), order.ToCommon(), offset, limit)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	result := make([]*Chemicals, len(res))
	for i, a := range res {
		result[i] = ChemicalsFromCommon(&a)
	}

	return result, nil
}

// HistoryOfActions is the resolver for the historyOfActions field.
func (r *queryResolver) HistoryOfActions(ctx context.Context, poolID common.ID, order Order, offset *int, limit *int) ([]*Action, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.actions.QueryActions(ctx, uuid.UUID(poolID), order.ToCommon(), offset, limit)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	result := make([]*Action, len(res))
	for i, a := range res {
		result[i], err = ActionFromCommon(a)
		if err != nil {
			return nil, castGQLError(ctx, err)
		}
	}

	return result, nil
}

// RecommendedChemicals is the resolver for the recommendedChemicals field.
func (r *queryResolver) RecommendedChemicals(ctx context.Context, poolID common.ID) ([]ChemicalValue, error) {
	if err := r.checkAccessToPool(ctx, uuid.UUID(poolID)); err != nil {
		return nil, err
	}

	res, err := r.estimator.RecommendedChemicals(ctx, uuid.UUID(poolID))
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return ChemicalValuesFromCommonProduct(res), nil
}

// MigrationStatus is the resolver for the migrationStatus field.
func (r *queryResolver) MigrationStatus(ctx context.Context, migrationID common.ID) (*Migration, error) {
	user, err := authPkg.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	status, err := r.migrator.Migration(ctx, user.ID, uuid.UUID(migrationID))
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	return MigrationFromCommon(status), nil
}

// OnCreatePool is the resolver for the onCreatePool field.
func (r *subscriptionResolver) OnCreatePool(ctx context.Context) (<-chan *Pool, error) {
	ch := make(chan *Pool)
	res, err := r.SubscribeOnCreate(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
			case pool := <-res:
				ch <- PoolFromCommon(pool)
			}
		}
	}()

	return ch, nil
}

// OnUpdatePool is the resolver for the onUpdatePool field.
func (r *subscriptionResolver) OnUpdatePool(ctx context.Context) (<-chan *Pool, error) {
	ch := make(chan *Pool)
	res, err := r.SubscribeOnUpdate(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
			case pool := <-res:
				ch <- PoolFromCommon(pool)
			}
		}
	}()

	return ch, nil
}

// OnDeletePool is the resolver for the onDeletePool field.
func (r *subscriptionResolver) OnDeletePool(ctx context.Context) (<-chan common.ID, error) {
	ch := make(chan common.ID)
	res, err := r.SubscribeOnDelete(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
			case id := <-res:
				ch <- common.ID(id)
			}
		}
	}()

	return ch, nil
}

// Pools is the resolver for the pools field.
func (r *userResolver) Pools(ctx context.Context, obj *User) ([]*Pool, error) {
	pools, err := r.poolData.List(ctx, obj.id)
	if err != nil {
		return nil, castGQLError(ctx, err)
	}

	result := make([]*Pool, len(pools))
	for i, pool := range pools {
		result[i] = PoolFromCommon(&pool)
	}

	return result, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Pool returns PoolResolver implementation.
func (r *Resolver) Pool() PoolResolver { return &poolResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type poolResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
