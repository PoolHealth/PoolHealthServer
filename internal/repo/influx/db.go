package influx

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type DB interface {
	measurement
	additives
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
