package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"e-commerce/pkg/models"
)

var httpClient = &http.Client{
    Timeout:  10 * time.Second,
}

func GetEstimatedAge (name string ) (*models.AgeEstimate, error){
    url := fmt.Sprintf("https://api.agify.io?name=%s", name)

    resp, err := httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result models.AgeEstimate
    if err := json.NewDecoder(resp.Body).Decode(&result) ; err != nil {
        return nil, err
    }
    return &result, nil
}