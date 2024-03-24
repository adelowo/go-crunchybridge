package gocrunchybridge

import (
	"context"
	"errors"
	"net/http"

	"github.com/adelowo/go-crunchybridge/internal/util"
)

type CertificateService service

type FetchTeamCertificateOption struct {
	TeamID string
}

func (c *CertificateService) Get(
	ctx context.Context,
	opts FetchTeamCertificateOption,
) (string, error) {
	if util.IsStringEmpty(opts.TeamID) {
		return "", errors.New("your team id is required")
	}

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return "", err
	}

	req, err := c.client.newRequest(http.MethodGet, "/teams/"+opts.TeamID+".pem", body)
	if err != nil {
		return "", err
	}

	var certificatePem string

	_, err = c.client.Do(ctx, req, &certificatePem)
	return certificatePem, err
}
