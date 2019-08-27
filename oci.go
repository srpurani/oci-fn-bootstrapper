package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
)
import "github.com/oracle/oci-go-sdk/core"

type fdkContext struct {
	Provider      string `json:"provider"`
	Registry      string `json:"registry"`
	ApiUrl        string `json:"api-url"`
	CallUrl       string `json:"call-url"`
	UserID        string `json:"oracle.user-id"`
	TenantID      string `json:"oracle.tenancy-id"`
	Fingerprint   string `json:"oracle.fingerprint"`
	PrivateKey    string `json:"oracle.key-file"`
	DisableCerts  bool   `json:"disable-certs"`
	CompartmentId string `json:"oracle.compartment-id"`
}

type ociState struct {
	VcnID    string `json:"vcn_id"`
	SubnetID string `json:"subnet_id"`
	TestApp  string `json:"test_app"`
	TestFunc string `json:"test_func"`
}

type Output struct {
	Status     string `json:"bootstrap_status"`
	ociState   `json:"oci_state"`
	fdkContext `json:"fdk_context"`
}

type RepoRequest struct {
	RepoDetails `contributesTo:"body"`
}

type RepoDetails struct {
	IsPublic *bool `mandatory:"true" json:"isPublic"`
}

func bootstrap(ctx context.Context, config *Config) (*Output, error) {
	decoded, err := base64.StdEncoding.DecodeString(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	configProvider := common.NewRawConfigurationProvider(config.TenantID, config.UserID, config.Region, config.Fingerprint, string(decoded), nil)
	vcnClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, wrap("failed to create vcnclient", err)
	}
	existingVcns, err := vcnClient.ListVcns(ctx, core.ListVcnsRequest{CompartmentId: common.String(config.CompartmentID)})
	if err != nil {
		return nil, wrap("failed to list vcns", err)
	}
	var vcnID string
	for _, existingVcn := range existingVcns.Items {
		if *existingVcn.DisplayName == config.VCNName {
			vcnID = *existingVcn.Id
			break
		}
	}
	// VCN not found, create one
	if vcnID == "" {
		createdVcn, err := vcnClient.CreateVcn(ctx, core.CreateVcnRequest{
			CreateVcnDetails: core.CreateVcnDetails{
				CompartmentId: common.String(config.CompartmentID),
				DisplayName:   common.String(config.VCNName),
				CidrBlock:     common.String("10.0.0.0/16"),
				DnsLabel:      common.String("testvcn"),
			},
		})
		if err != nil {
			return nil, wrap("failed to create vcn", err)
		}
		vcnID = *createdVcn.Id
	}
	existingSubnets, err := vcnClient.ListSubnets(ctx, core.ListSubnetsRequest{
		CompartmentId: common.String(config.CompartmentID),
		VcnId:         common.String(vcnID),
	})
	if err != nil {
		return nil, wrap("failed to list subnets", err)
	}
	var subnetID string
	var routeTableID string
	for _, existingSubnet := range existingSubnets.Items {
		if *existingSubnet.DisplayName == config.SubnetRegional {
			subnetID = *existingSubnet.Id
			routeTableID = *existingSubnet.RouteTableId
			break
		}
	}
	// Subnet does not exist, create one
	if subnetID == "" {
		createdSubnet, err := vcnClient.CreateSubnet(ctx, core.CreateSubnetRequest{
			CreateSubnetDetails: core.CreateSubnetDetails{
				DisplayName:   common.String(config.SubnetRegional),
				CompartmentId: common.String(config.CompartmentID),
				VcnId:         common.String(vcnID),
				CidrBlock:     common.String("10.0.1.0/24"),
			},
		})
		if err != nil {
			return nil, wrap("failed to create subnet", err)
		}
		subnetID = *createdSubnet.Id
		routeTableID = *createdSubnet.RouteTableId
	}
	// Create Internet Gateway if it does not exist
	existingInternetGateways, err := vcnClient.ListInternetGateways(ctx, core.ListInternetGatewaysRequest{
		VcnId:         common.String(vcnID),
		CompartmentId: common.String(config.CompartmentID),
	})
	if err != nil {
		return nil, wrap("failed to list internet gateways", err)
	}
	var igID string
	for _, existingIG := range existingInternetGateways.Items {
		igID = *existingIG.Id
		break
	}
	if igID == "" {
		// IG does not exist, create one
		createdIG, err := vcnClient.CreateInternetGateway(ctx, core.CreateInternetGatewayRequest{
			CreateInternetGatewayDetails: core.CreateInternetGatewayDetails{
				CompartmentId: common.String(config.CompartmentID),
				DisplayName:   common.String("ig"),
				VcnId:         common.String(vcnID),
				IsEnabled:     common.Bool(true),
			},
		})
		if err != nil {
			return nil, wrap("failed to create internet gateway", err)
		}
		igID = *createdIG.Id
	}
	// Add a routing rule to internet gateway, if needed
	routeTable, err := vcnClient.GetRouteTable(ctx, core.GetRouteTableRequest{
		RtId: common.String(routeTableID),
	})
	if err != nil {
		return nil, wrap("failed to get route table", err)
	}
	var routingRuleFound bool
	for _, rule := range routeTable.RouteRules {
		if *rule.NetworkEntityId == igID {
			routingRuleFound = true
			break
		}
	}
	if !routingRuleFound {
		_, err := vcnClient.UpdateRouteTable(ctx, core.UpdateRouteTableRequest{
			RtId: common.String(routeTableID),
			UpdateRouteTableDetails: core.UpdateRouteTableDetails{
				RouteRules: []core.RouteRule{core.RouteRule{
					NetworkEntityId: common.String(igID),
					Destination:     common.String("0.0.0.0/0"),
					DestinationType: core.RouteRuleDestinationTypeCidrBlock,
				}},
			},
		})
		if err != nil {
			return nil, wrap("failed to update route table", err)
		}
	}
	// Try to create repository
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return nil, err
	}
	baseClient.Host = "https://iad.ocir.io"
	request := common.MakeDefaultHTTPRequest("GET", fmt.Sprintf("/20180419/docker/repos/%s/%s", config.TenantName, config.RepositoryName))
	_, err = baseClient.Call(ctx, &request)
	if _, ok := err.(common.ServiceError); ok {
		request, err := common.MakeDefaultHTTPRequestWithTaggedStruct("POST", fmt.Sprintf("/20180419/docker/repos/%s/%s", config.TenantName, config.RepositoryName), RepoRequest{
			RepoDetails: RepoDetails{
				IsPublic: common.Bool(true),
			},
		})
		if err != nil {
			return nil, wrap("failed to create repo create request", err)
		}
		_, err = baseClient.Call(ctx, &request)
		if err != nil {
			return nil, wrap("failed to create repo", err)
		}
	}
	output := &Output{
		Status: "succeeded",
		ociState: ociState{
			SubnetID: subnetID,
			VcnID:    vcnID,
			TestApp:  "test-app",
			TestFunc: "test-func",
		},
		fdkContext: fdkContext{
			CompartmentId: config.CompartmentID,
			//PrivateKey:    config.PrivateKey,
			Fingerprint:  config.Fingerprint,
			UserID:       config.UserID,
			TenantID:     config.TenantID,
			ApiUrl:       fmt.Sprintf("https://functions.%s.oraclecloud.com", config.Region),
			DisableCerts: true,
			Provider:     "oracle",
			Registry:     fmt.Sprintf("iad.ocir.io/%s/%s", config.TenantName, config.RepositoryName),
		},
	}

	return output, nil
}

func wrap(msg string, err error) error {
	return fmt.Errorf("msg: %s, err: %s", msg, err)
}
