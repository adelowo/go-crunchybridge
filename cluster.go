package gocrunchybridge

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adelowo/go-crunchybridge/internal/util"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/google/go-querystring/query"
)

type ClusterService service

// ENUM(production)
type ClusterEnvironment string

// ENUM(aws,gcp,azure)
type ClusterProvider string

func (c *ClusterService) Create(ctx context.Context,
	opts *CreateClusterOptions,
) (Cluster, error) {
	var cluster Cluster

	if err := opts.Validate(); err != nil {
		return cluster, err
	}

	body, err := ToReader(opts)
	if err != nil {
		return cluster, err
	}

	req, err := c.client.newRequest(http.MethodPost, "/clusters", body)
	if err != nil {
		return cluster, err
	}

	_, err = c.client.Do(ctx, req, &cluster)
	return cluster, err
}

func (c *ClusterService) Delete(ctx context.Context,
	opts FetchClusterOptions) error {

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	if err := opts.ID.IsValid(); err != nil {
		return fmt.Errorf("cluster_id: %w", err)
	}

	req, err := c.client.newRequest(http.MethodDelete,
		fmt.Sprintf("/clusters/%s", opts.ID), body)
	if err != nil {
		return err
	}

	_, err = c.client.Do(ctx, req, nil)
	return err
}

func (c *ClusterService) Ping(ctx context.Context,
	opts FetchClusterOptions) error {

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	if err := opts.ID.IsValid(); err != nil {
		return fmt.Errorf("cluster_id: %w", err)
	}

	req, err := c.client.newRequest(http.MethodPut,
		fmt.Sprintf("/clusters/%s/actions/ping", opts.ID), body)
	if err != nil {
		return err
	}

	_, err = c.client.Do(ctx, req, nil)
	return err
}

func (c *ClusterService) List(ctx context.Context,
	opts ListClusterOptions) (ListClusterResponse, error) {

	var clusters ListClusterResponse

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return clusters, err
	}

	v, err := query.Values(opts)
	if err != nil {
		return clusters, err
	}

	endpoint := "/clusters?" + v.Encode()

	req, err := c.client.newRequest(http.MethodGet, endpoint, body)
	if err != nil {
		return clusters, err
	}

	_, err = c.client.Do(ctx, req, &clusters)
	return clusters, err
}

type FetchClusterOptions struct {
	ID EID
}

// ENUM(asc,desc)
type ListFilterOrderType string

// ENUM(id,name)
type ListOrderFieldType string

type ListClusterOptions struct {
	Order      ListFilterOrderType `url:"order,omitempty"`
	OrderField string              `url:"order_field,omitempty"`
	TeamID     EID                 `url:"team_id,omitempty"`
	Limit      int                 `url:"limit,omitempty"`
}

type ListClusterResponse struct {
	HasMore    bool      `json:"has_more,omitempty"`
	Clusters   []Cluster `json:"clusters,omitempty"`
	NextCursor EID       `json:"next_cursor,omitempty"`
}

type CreateClusterOptions struct {
	// Hunan readable name for the cluster. If non is provided,
	// a default name would be generated for the cluster
	Name   string `json:"name,omitempty"`
	PlanID string `json:"plan_id,omitempty"`
	// You must have admin access to create a cluster
	TeamID EID `json:"team_id,omitempty"`

	ClusterGroupID EID `json:"cluster_group_id,omitempty"`

	HighlyAvailable   bool `json:"is_ha,omitempty"`
	KeychainID        EID  `json:"keychain_id,omitempty"`
	NetworkID         EID  `json:"network_id,omitempty"`
	PostgresVersionID int  `json:"postgres_version_id,omitempty"`
	// You cannot set this and NetworkID together
	RegionID    string          `json:"region_id,omitempty"`
	StorageSize int             `json:"storage,omitempty"`
	ProviderID  ClusterProvider `json:"provider_id,omitempty"`

	Environment ClusterEnvironment `json:"environment,omitempty"`
}

func (c *CreateClusterOptions) Validate() error {
	if err := c.TeamID.IsValid(); err != nil {
		return fmt.Errorf("team id: %v", err)
	}

	if util.IsStringEmpty(c.Name) {
		c.Name = petname.Generate(3, "-")
	}

	if util.IsStringEmpty(c.PlanID) {
		return errors.New("please provide a plan ID")
	}

	splitted := strings.Split(c.PlanID, "-")
	if len(splitted) != 2 {
		return errors.New("please provide a valid plan ID. Plan ID must be in this format: hobby-1,standard-1")
	}

	if !util.IsStringEmpty(c.ClusterGroupID.String()) {
		if err := c.ClusterGroupID.IsValid(); err != nil {
			return fmt.Errorf("cluster_group_id: %v", err)
		}
	}

	if !c.Environment.IsValid() {
		return errors.New("environment: please provide a valid environment value")
	}

	if !util.IsStringEmpty(c.KeychainID.String()) {
		if err := c.KeychainID.IsValid(); err != nil {
			return fmt.Errorf("keychain_id: %v", err)
		}
	}

	if !util.IsStringEmpty(c.NetworkID.String()) {
		if err := c.NetworkID.IsValid(); err != nil {
			return fmt.Errorf("network_id: %v", err)
		}
	}

	if c.PostgresVersionID < 12 {
		return errors.New("minimum supported postgres version is 12")
	}

	if !c.ProviderID.IsValid() {
		return errors.New("provider_id: please provide a valid and supported cloud provider")
	}

	if util.IsStringEmpty(c.RegionID) {
		return errors.New("region_id: please provide a region ID")
	}

	if c.StorageSize < 5 {
		return errors.New("minimum storage is 5GB")
	}

	return nil
}

type ClusterDiskUsage struct {
	DiskAvailableMb int `json:"disk_available_mb"`
	DiskTotalSizeMb int `json:"disk_total_size_mb"`
	DiskUsedMb      int `json:"disk_used_mb"`
}

type Cluster struct {
	ClusterID              any                `json:"cluster_id"`
	CPU                    int                `json:"cpu"`
	CreatedAt              time.Time          `json:"created_at"`
	DiskUsage              ClusterDiskUsage   `json:"disk_usage"`
	Environment            ClusterEnvironment `json:"environment"`
	Host                   string             `json:"host"`
	ID                     EID                `json:"id"`
	IsHighlyAvailable      bool               `json:"is_ha"`
	IsProtected            bool               `json:"is_protected"`
	IsSuspended            bool               `json:"is_suspended"`
	MaintenanceWindowStart int                `json:"maintenance_window_start"`
	MajorVersion           int                `json:"major_version"`
	Memory                 float64            `json:"memory"`
	Name                   string             `json:"name"`
	NetworkID              EID                `json:"network_id"`
	ParentID               EID                `json:"parent_id"`
	PlanID                 string             `json:"plan_id"`
	PostgresVersionID      string             `json:"postgres_version_id"`
	ProviderID             EID                `json:"provider_id"`
	RegionID               EID                `json:"region_id"`
	Replicas               []struct {
		ClusterID              any                `json:"cluster_id"`
		CPU                    int                `json:"cpu"`
		CreatedAt              time.Time          `json:"created_at"`
		DiskUsage              ClusterDiskUsage   `json:"disk_usage"`
		Environment            ClusterEnvironment `json:"environment"`
		Host                   string             `json:"host"`
		ID                     EID                `json:"id"`
		IsHighlyAvailable      bool               `json:"is_ha"`
		IsProtected            bool               `json:"is_protected"`
		IsSuspended            bool               `json:"is_suspended"`
		MaintenanceWindowStart int                `json:"maintenance_window_start"`
		MajorVersion           int                `json:"major_version"`
		Memory                 float64            `json:"memory"`
		Name                   string             `json:"name"`
		NetworkID              EID                `json:"network_id"`
		ParentID               EID                `json:"parent_id"`
		PlanID                 string             `json:"plan_id"`
		PostgresVersionID      string             `json:"postgres_version_id"`
		ProviderID             EID                `json:"provider_id"`
		RegionID               EID                `json:"region_id"`
		ResetStatsWeekly       bool               `json:"reset_stats_weekly"`
		Storage                int                `json:"storage"`
		TailscaleActive        bool               `json:"tailscale_active"`
		TeamID                 EID                `json:"team_id"`
		UpdatedAt              time.Time          `json:"updated_at"`
	} `json:"replicas"`
	ResetStatsWeekly bool      `json:"reset_stats_weekly"`
	Storage          int       `json:"storage"`
	TailscaleActive  bool      `json:"tailscale_active"`
	TeamID           EID       `json:"team_id"`
	UpdatedAt        time.Time `json:"updated_at"`
}
