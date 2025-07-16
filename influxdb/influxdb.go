package influxdb

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Config struct {
	Url     string
	Token   string
	Org     string
	Bucket  string
	Timeout int
}

type Client struct {
	client influxdb2.Client
	org    string
	bucket string
}

func (cfg Config) NewClient() *Client {
	client := influxdb2.NewClient(cfg.Url, cfg.Token)
	return &Client{
		client: client,
		org:    cfg.Org,
		bucket: cfg.Bucket,
	}
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) WritePoint(point *write.Point) error {
	writeAPI := c.client.WriteAPIBlocking(c.org, c.bucket)
	return writeAPI.WritePoint(context.Background(), point)
}

func (c *Client) BatchWritePoints(points []*write.Point) error {
	writeAPI := c.client.WriteAPIBlocking(c.org, c.bucket)
	return writeAPI.WritePoint(context.Background(), points...)
}

func (c *Client) Query(fluxQuery string) ([]map[string]interface{}, error) {
	queryAPI := c.client.QueryAPI(c.org)
	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var results []map[string]interface{}
	for result.Next() {
		results = append(results, result.Record().Values())
	}

	if result.Err() != nil {
		return results, result.Err()
	}

	return results, nil
}
