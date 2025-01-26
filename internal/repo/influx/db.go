package influx

import (
	"context"
	"errors"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/http"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type DB interface {
	measurement
	additives
	action
}

type db struct {
	client    influxdb2.Client
	writeAPI  api.WriteAPIBlocking
	deleteAPI api.DeleteAPI
	queryAPI  api.QueryAPI
	bucket    string
	org       string

	log log.Logger
}

func New(addr, token, org, bucket string, logger log.Logger) (DB, error) {
	client := influxdb2.NewClient(addr, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)
	deleteAPI := client.DeleteAPI()
	queryAPI := client.QueryAPI(org)

	return &db{
		client:    client,
		writeAPI:  writeAPI,
		queryAPI:  queryAPI,
		deleteAPI: deleteAPI,
		org:       org,
		bucket:    bucket,
		log:       logger,
	}, nil
}

func (d *db) writePoint(ctx context.Context, point ...*write.Point) error {
	err := d.writeAPI.WritePoint(ctx, point...)
	if err == nil {
		return nil
	}

	httpErr := &http.Error{}
	if !errors.As(err, &httpErr) {
		return err
	}

	if httpErr.StatusCode == 422 {
		return common.ErrRecordAlreadyExists
	}

	return err
}
