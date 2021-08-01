package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

type InfluxDb struct {
	DBclient api.WriteAPI
}

func NewInfluxDb() *InfluxDb {
	client := influxdb2.NewClientWithOptions("http://localhost:8086", "Jggte8lIdnjPAyxPi4njs-DhSOjuuYm2R4V4E7STq1xSIVv1TGKemdLE5fs_UpoUxCdlxDsTOL7RJ4mTEKGzZA==",
		influxdb2.DefaultOptions().SetBatchSize(20))

	writeAPI := client.WriteAPI("myorg", "mybucket")
	return &InfluxDb{
		DBclient: writeAPI,
	}
}
