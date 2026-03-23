package models

type AgeEstimate struct {
    Name string `json:"name", cql:"name"`
    Age int `json:"age", cql:"age"`
    Count int `json:"count", cql:"count"`
}