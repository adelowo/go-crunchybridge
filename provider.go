package gocrunchybridge

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
)

type ProviderService service

type FetchProviderOptions struct {
	Provider ClusterProvider
}

type ProviderResponseArray struct {
	Providers []struct {
		Disk struct {
			Rate int `json:"rate"`
		} `json:"disk"`
		DisplayName string `json:"display_name"`
		ID          string `json:"id"`
		Plans       []struct {
			CPU                       int     `json:"cpu"`
			DisplayName               string  `json:"display_name"`
			ID                        string  `json:"id"`
			IopsBaseline              int     `json:"iops_baseline"`
			IopsMaximum               int     `json:"iops_maximum"`
			MaximumPerformanceLimited bool    `json:"maximum_performance_limited"`
			Memory                    float64 `json:"memory"`
			Rate                      int     `json:"rate"`
		} `json:"plans"`
		Regions []struct {
			DisplayName string  `json:"display_name"`
			ID          string  `json:"id"`
			Location    string  `json:"location"`
			Multiplier  float64 `json:"multiplier"`
		} `json:"regions"`
	} `json:"providers"`
}

func (p *ProviderService) List(ctx context.Context,
	opts FetchProviderOptions) (ProviderResponse, error) {

	var resp ProviderResponse

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := p.client.newRequest(http.MethodGet, "/providers", body)
	if err != nil {
		return resp, err
	}

	var r interface{}
	_, err = p.client.Do(ctx, req, &r)

	json.NewEncoder(os.Stdout).Encode(r)

	return resp, err
}
