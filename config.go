package main

import (
	"encoding/json"
	"io"
)

type Config struct {
	UserID string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Region string `json:"region"`
	Fingerprint string `json:"fingerprint"`
	PrivateKey  string `json:"private_key"`
	TenantName string `json:"tenant_name"`

	CompartmentID string `json:"compartment_id"`
	VCNName       string  `json:"vcn_name"`
	SubnetRegional string `json:"regional_subnet"`
	RepositoryName string `json:"repo_name"`
}

func readConfig(in io.Reader) (*Config, error){
	c := &Config{}
	if err := json.NewDecoder(in).Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}