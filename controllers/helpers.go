package controllers

import (
	"bytes"
	"encoding/json"

	"strconv"
	"time"

	"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
)

func isValidSKU(skuID uint) bool {
	skuIDStr := strconv.Itoa(int(skuID))
	config := interservice_client.Config{
		ServiceName: "user-service",
		BaseURL:     "http://localhost:8081/api/V1/skus/",
		Timeout:     5 * time.Second,
	}
	client, err := interservice_client.NewClientWithConfig(config)
	if err != nil {
		return false
	}
	url := config.BaseURL + "validate/" + skuIDStr
	body := map[string]string{
		"hub_id": "",
		"skus":   "",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return false
	}
	req := &http.Request{
		Url:     url, // Use configured URL
		Body:    bytes.NewReader(bodyBytes),
		Timeout: 7 * time.Second,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}
	resp, _ := client.Get(req, "/")
	if resp == nil {
		return false
	}
	return resp.StatusCode() == 200
}

func isValidHub(hubID uint) bool {
	hubIDStr := strconv.Itoa(int(hubID))
	config := interservice_client.Config{
		ServiceName: "user-service",
		BaseURL:     "http://localhost:8081/api/V1/hubs/",
		Timeout:     5 * time.Second,
	}
	client, err := interservice_client.NewClientWithConfig(config)
	if err != nil {
		return false
	}
	url := config.BaseURL + "validate/" + hubIDStr
	body := map[string]string{
		"hub_id": "",
		"skus":   "",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return false
	}
	req := &http.Request{
		Url:     url,
		Body:    bytes.NewReader(bodyBytes),
		Timeout: 7 * time.Second,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}
	resp, _ := client.Get(req, "")
	if resp == nil {
		return false
	}
	return resp.StatusCode() == 200

}
