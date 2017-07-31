package influxdbinterface

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	influxdbAddress = "http://kless-influxdb.kless:8086"
	db              = "klessdb"
	seriesName      = "kless"
)

// InfluxdbInterface enables interaction with an InfluxDB instance
type InfluxdbInterface struct {
}

//KlessEvent defines what gets stored in InfluxDB
type KlessEvent struct {
	Timestamp    string
	Namespace    string
	EventHandler string
	Version      string
	PodName      string
	RequestSize  int64
	ResponseSize int64
	ResponseTime int64
}

// AddEvent adds an event to InfluxDB
func (i *InfluxdbInterface) AddEvent(e *KlessEvent) error {
	fmt.Println("Entering AddEvent")

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxdbAddress,
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	fmt.Println("Connected")

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "ms",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// Create a point and add to batch
	tags := map[string]string{"namespace": e.Namespace, "handler": e.EventHandler, "version": e.Version, "podname": e.PodName}
	fields := map[string]interface{}{
		"reqSize":  e.RequestSize,
		"respSize": e.ResponseSize,
		"respTime": e.ResponseTime,
	}
	pt, err := client.NewPoint(seriesName, tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)

	fmt.Println("Leaving AddEvent")

	return nil
}

//GetEvents retrieves the list of events from InfluxDB, this currently retrieves all events which obviously needs to be fixed
func (i *InfluxdbInterface) GetEvents() (events []KlessEvent, e error) {

	fmt.Println("Entering GetEvents")

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxdbAddress,
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	fmt.Println("Connected")

	q := fmt.Sprintf("SELECT time, namespace, handler, version, podname, reqSize, respSize, respTime FROM %s", seriesName)

	res, err := queryDB(c, q)
	if err != nil {
		log.Fatal(err)
	}

	if len(res) > 0 && len(res[0].Series) > 0 {
		for _, row := range res[0].Series[0].Values {
			t, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Fatal(err)
			}
			requestSize, err := row[5].(json.Number).Int64()
			if err != nil {
				log.Fatal(err)
			}
			responseSize, err := row[6].(json.Number).Int64()
			if err != nil {
				log.Fatal(err)
			}
			responseTime, err := row[7].(json.Number).Int64()
			if err != nil {
				log.Fatal(err)
			}
			event := KlessEvent{
				Timestamp:    t.Format(time.RFC3339Nano),
				Namespace:    row[1].(string),
				EventHandler: row[2].(string),
				Version:      row[3].(string),
				PodName:      row[4].(string),
				RequestSize:  requestSize,
				ResponseSize: responseSize,
				ResponseTime: responseTime,
			}
			events = append(events, event)
		}
	}

	fmt.Println("Leaving GetEvents")

	return events, nil
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
