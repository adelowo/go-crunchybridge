package gocrunchybridge

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ProviderService service

type FetchProviderOptions struct {
	Provider ClusterProvider
}

type ProviderResponse struct {
	Providers []ProviderPlan `json:"providers"`
}

type ProviderPlan struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Plans       []struct {
		CPU                       int     `json:"cpu"`
		DisplayName               string  `json:"display_name"`
		ID                        EID     `json:"id"`
		IopsBaseline              int     `json:"iops_baseline"`
		IopsMaximum               int     `json:"iops_maximum"`
		MaximumPerformanceLimited bool    `json:"maximum_performance_limited"`
		Memory                    float64 `json:"memory"`
		Rate                      int     `json:"rate"`
	} `json:"plans"`
	Regions []struct {
		DisplayName string `json:"display_name"`
		ID          EID    `json:"id"`
		Location    string `json:"location"`
	} `json:"regions"`
}

func (p *ProviderService) Get(ctx context.Context,
	opts FetchProviderOptions) (ProviderPlan, error) {

	var resp ProviderPlan

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := p.client.newRequest(http.MethodGet, "/providers", body)
	if err != nil {
		return resp, err
	}

	var plan ProviderResponse
	_, err = p.client.Do(ctx, req, &plan)
	if err != nil {
		return resp, err
	}

	for _, provider := range plan.Providers {
		if strings.ToLower(provider.ID) == strings.ToLower(opts.Provider.String()) {
			return provider, nil
		}
	}

	return resp, fmt.Errorf("provider (%s) does not exists", opts.Provider.String())
}
