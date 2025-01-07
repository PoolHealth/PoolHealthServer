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
	additivesTable = "additives"
)

type additives interface {
	CreateAdditive(ctx context.Context, rec *common.Additives) error
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Additives, error)
}

func (d *db) CreateAdditive(ctx context.Context, rec *common.Additives) error {
	p := influxdb2.NewPointWithMeasurement(additivesTable).
		AddTag(poolIDKey, rec.PoolID.String()).
		SetTime(rec.CreatedAt)

	for k, v := range rec.Products {
		p.AddField(common.ChemicalProductNames[k], v)
	}

	return d.writeAPI.WritePoint(ctx, p)
}

func (d *db) QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Additives, error) {
	query := fmt.Sprintf(`from(bucket:"%s") 
	|> range(start: -1y) 
	|> filter(fn: (r) => r._measurement == "%s" and r.poolID == "%s")
	|> window(every: 1d)
	|> sum()`, d.bucket, additivesTable, poolID.String())

	d.log.Info(query)
	result, err := d.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	data := map[time.Time]map[common.ChemicalProduct]float64{}
	times := make([]time.Time, 0)
	for result.Next() {
		record := result.Record()
		t := record.Start()
		if _, ok := data[t]; !ok {
			times = append(times, record.Start())
			data[t] = map[common.ChemicalProduct]float64{}
		}

		product, ok := common.ChemicalProductNamesToChemicalProduct[record.Field()]
		if ok {
			data[t][product] = record.Value().(float64)
		}
	}

	if order == common.OrderAsc {
		slices.SortFunc(times, func(a, b time.Time) int { return a.Compare(b) })
	} else {
		slices.SortFunc(times, func(a, b time.Time) int { return b.Compare(a) })
	}

	res := make([]common.Additives, len(times))
	for i, t := range times {
		res[i] = common.Additives{
			PoolID:    poolID,
			CreatedAt: t,
			Products:  data[t],
		}
	}

	return res, nil
}
