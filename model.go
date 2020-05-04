package main

import "time"

//Metrics ...
type Metrics struct {
	ServiceName string     `json:"ServiceName"`
	MetricValue float64    `json:"MetricValue"`
	MetricName  string     `json:"MetricName"`
	Time        *time.Time `json:"Time"`
	Status      int        `json:"Status"`
}
