package gocrunchybridge

type ClusterService service

type CreateClusterOptions struct {
	// Hunan readable name for the cluster. If non is provided,
	// a default name would be generated for the cluster
	Name   string
	PlanID string
}
