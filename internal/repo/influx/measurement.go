package influx

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/PoolHealth/PoolHealthServer/common"
)

const (
	measurementTable = "measurement"
	poolIDKey        = "poolID"
	chlorineKey      = "chlorine"
	phKey            = "ph"
	alkalinityKey    = "alkalinity"
)

type measurement interface {
	CreateMeasurement(ctx context.Context, rec common.Measurement) error
	QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Measurement, error)
	QueryLastMeasurement(ctx context.Context, poolID uuid.UUID) ([]common.Measurement, error)
}

func (d *db) CreateMeasurement(ctx context.Context, rec common.Measurement) error {
	p := influxdb2.NewPointWithMeasurement(measurementTable).
		AddTag(poolIDKey, rec.PoolID.String()).
		SetTime(rec.CreatedAt)

	if rec.Chlorine.Valid {
		p.AddField(chlorineKey, rec.Chlorine.Float64)
	}

	if rec.PH.Valid {
		p.AddField(phKey, rec.PH.Float64)
	}

	if rec.Alkalinity.Valid {
		p.AddField(alkalinityKey, rec.Alkalinity.Float64)
	}

	return d.writeAPI.WritePoint(ctx, p)
}

func (d *db) QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Measurement, error) {
	query := fmt.Sprintf(`from(bucket:"%s") 
|> range(start: -1y) 
|> filter(fn: (r) => r._measurement == "%s" and r.poolID == "%s")`, d.bucket, measurementTable, poolID.String())

	return d.queryMeasurement(ctx, query, poolID, order)
}

func (d *db) QueryLastMeasurement(ctx context.Context, poolID uuid.UUID) ([]common.Measurement, error) {
	query := fmt.Sprintf(`from(bucket:"%s")
	|> range(start: -1y)
	|> filter(fn: (r) => r._measurement == "%s" and r.poolID == "%s")
	|> window(every: 1d)
	|> last()`, d.bucket, measurementTable, poolID.String())

	return d.queryMeasurement(ctx, query, poolID, common.OrderDesc)
}

func (d *db) queryMeasurement(ctx context.Context, query string, poolID uuid.UUID, order common.Order) ([]common.Measurement, error) {
	d.log.Info(query)
	result, err := d.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	data := map[time.Time]map[string]null.Float{}

	times := make([]time.Time, 0)

	// Use Next() to iterate over query result lines
	for result.Next() {
		// Observe when there is new grouping key producing new table
		if result.TableChanged() {
			d.log.WithField("table", result.TableMetadata().String()).Info("changed")
		}

		d.log.
			WithField("time", result.Record().Time()).
			WithField("value", result.Record().Value()).
			WithField("field", result.Record().Field()).
			Info("record")
		// read result
		record := result.Record()

		el, ok := data[record.Time()]
		if !ok {
			times = append(times, record.Time())
			data[record.Time()] = map[string]null.Float{record.Field(): null.FloatFrom(record.Value().(float64))}
			continue
		}

		el[record.Field()] = null.FloatFrom(record.Value().(float64))
		data[record.Time()] = el
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	if order == common.OrderAsc {
		slices.SortFunc(times, func(a, b time.Time) int { return a.Compare(b) })
	} else {
		slices.SortFunc(times, func(a, b time.Time) int { return b.Compare(a) })
	}

	res := make([]common.Measurement, len(data))

	for i, t := range times {
		el := common.Measurement{
			CreatedAt: t,
			PoolID:    poolID,
		}

		if val, ok := data[t][chlorineKey]; ok {
			el.Chlorine = val
		}

		if val, ok := data[t][phKey]; ok {
			el.PH = val
		}

		if val, ok := data[t][alkalinityKey]; ok {
			el.Alkalinity = val
		}

		res[i] = el
	}

	return res, nil
}
