package main

type action struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}