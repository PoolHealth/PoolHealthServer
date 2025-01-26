package influx

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/PoolHealth/PoolHealthServer/common"
)

const (
	actionsTable = "actions"
)

type action interface {
	LogActions(ctx context.Context, poolID uuid.UUID, actions *common.Action) error
	QueryActions(
		ctx context.Context,
		poolID uuid.UUID,
		order common.Order,
	) ([]common.Action, error)
	DeleteAction(ctx context.Context, poolID uuid.UUID, createdAt time.Time) error
}

func (d *db) LogActions(ctx context.Context, poolID uuid.UUID, actions *common.Action) error {
	p := influxdb2.NewPointWithMeasurement(actionsTable).
		AddTag(poolIDKey, poolID.String()).
		SetTime(actions.CreatedAt)

	for _, v := range actions.Types {
		p.AddField(common.ActionTypeNames[v], true)
	}

	return d.writePoint(ctx, p)
}

func (d *db) QueryActions(
	ctx context.Context,
	poolID uuid.UUID,
	order common.Order,
) ([]common.Action, error) {
	query := fmt.Sprintf(`from(bucket:"%s") 
	|> range(start: -1y) 
	|> filter(fn: (r) => r._measurement == "%s" and r.poolID == "%s")
	|> window(every: 1m)
	|> count()`, d.bucket, actionsTable, poolID.String())

	return d.queryActions(ctx, query, order)
}

func (d *db) DeleteAction(ctx context.Context, poolID uuid.UUID, createdAt time.Time) error {
	return d.deleteAPI.DeleteWithName(
		ctx,
		d.org,
		d.bucket,
		createdAt,
		createdAt.Add(time.Minute),
		fmt.Sprintf(`_measurement="%s" and poolID="%s"`, actionsTable, poolID.String()),
	)
}

func (d *db) queryActions(ctx context.Context, query string, order common.Order) ([]common.Action, error) {
	d.log.Info(query)

	result, err := d.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	data := map[time.Time][]common.ActionType{}
	times := make([]time.Time, 0)

	for result.Next() {
		record := result.Record()
		t := record.Start()

		if record.ValueByKey("_time") != nil {
			t = record.Time()
		}

		if _, ok := data[t]; !ok {
			times = append(times, t)
			data[t] = []common.ActionType{}
		}

		actionType, ok := common.ActionTypeNamesToActionType[record.Field()]
		if ok {
			data[t] = append(data[t], actionType)
		}
	}

	if order == common.OrderAsc {
		slices.SortFunc(times, func(a, b time.Time) int { return a.Compare(b) })
	} else {
		slices.SortFunc(times, func(a, b time.Time) int { return b.Compare(a) })
	}

	res := make([]common.Action, len(times))
	for i, t := range times {
		res[i] = common.Action{
			CreatedAt: t,
			Types:     data[t],
		}
	}

	return res, nil
}
