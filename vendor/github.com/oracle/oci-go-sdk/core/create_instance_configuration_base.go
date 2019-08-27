// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// CreateInstanceConfigurationBase Creation details for an instance configuration.
type CreateInstanceConfigurationBase interface {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment
	// containing the instance configuration.
	GetCompartmentId() *string

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	GetDefinedTags() map[string]map[string]interface{}

	// A user-friendly name for the instance configuration.  Does not have to be unique,
	// and it's changeable. Avoid entering confidential information.
	GetDisplayName() *string

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	GetFreeformTags() map[string]string
}

type createinstanceconfigurationbase struct {
	JsonData      []byte
	CompartmentId *string                           `mandatory:"true" json:"compartmentId"`
	DefinedTags   map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
	DisplayName   *string                           `mandatory:"false" json:"displayName"`
	FreeformTags  map[string]string                 `mandatory:"false" json:"freeformTags"`
	Source        string                            `json:"source"`
}

// UnmarshalJSON unmarshals json
func (m *createinstanceconfigurationbase) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalercreateinstanceconfigurationbase createinstanceconfigurationbase
	s := struct {
		Model Unmarshalercreateinstanceconfigurationbase
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.CompartmentId = s.Model.CompartmentId
	m.DefinedTags = s.Model.DefinedTags
	m.DisplayName = s.Model.DisplayName
	m.FreeformTags = s.Model.FreeformTags
	m.Source = s.Model.Source

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *createinstanceconfigurationbase) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Source {
	case "NONE":
		mm := CreateInstanceConfigurationDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "INSTANCE":
		mm := CreateInstanceConfigurationFromInstanceDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetCompartmentId returns CompartmentId
func (m createinstanceconfigurationbase) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetDefinedTags returns DefinedTags
func (m createinstanceconfigurationbase) GetDefinedTags() map[string]map[string]interface{} {
	return m.DefinedTags
}

//GetDisplayName returns DisplayName
func (m createinstanceconfigurationbase) GetDisplayName() *string {
	return m.DisplayName
}

//GetFreeformTags returns FreeformTags
func (m createinstanceconfigurationbase) GetFreeformTags() map[string]string {
	return m.FreeformTags
}

func (m createinstanceconfigurationbase) String() string {
	return common.PointerString(m)
}
