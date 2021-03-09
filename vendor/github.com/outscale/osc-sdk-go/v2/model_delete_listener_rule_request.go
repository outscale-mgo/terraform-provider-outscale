/*
 * 3DS OUTSCALE API
 *
 * Welcome to the 3DS OUTSCALE's API documentation.<br /><br />  The 3DS OUTSCALE API enables you to manage your resources in the 3DS OUTSCALE Cloud. This documentation describes the different actions available along with code examples.<br /><br />  Note that the 3DS OUTSCALE Cloud is compatible with Amazon Web Services (AWS) APIs, but some resources have different names in AWS than in the 3DS OUTSCALE API. You can find a list of the differences [here](https://wiki.outscale.net/display/EN/3DS+OUTSCALE+APIs+Reference).<br /><br />  You can also manage your resources using the [Cockpit](https://wiki.outscale.net/display/EN/About+Cockpit) web interface.
 *
 * API version: 1.7
 * Contact: support@outscale.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package osc

import (
	"encoding/json"
)

// DeleteListenerRuleRequest struct for DeleteListenerRuleRequest
type DeleteListenerRuleRequest struct {
	// If true, checks whether you have the required permissions to perform the action.
	DryRun *bool `json:"DryRun,omitempty"`
	// The name of the rule you want to delete.
	ListenerRuleName string `json:"ListenerRuleName"`
}

// NewDeleteListenerRuleRequest instantiates a new DeleteListenerRuleRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDeleteListenerRuleRequest(listenerRuleName string) *DeleteListenerRuleRequest {
	this := DeleteListenerRuleRequest{}
	this.ListenerRuleName = listenerRuleName
	return &this
}

// NewDeleteListenerRuleRequestWithDefaults instantiates a new DeleteListenerRuleRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDeleteListenerRuleRequestWithDefaults() *DeleteListenerRuleRequest {
	this := DeleteListenerRuleRequest{}
	return &this
}

// GetDryRun returns the DryRun field value if set, zero value otherwise.
func (o *DeleteListenerRuleRequest) GetDryRun() bool {
	if o == nil || o.DryRun == nil {
		var ret bool
		return ret
	}
	return *o.DryRun
}

// GetDryRunOk returns a tuple with the DryRun field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DeleteListenerRuleRequest) GetDryRunOk() (*bool, bool) {
	if o == nil || o.DryRun == nil {
		return nil, false
	}
	return o.DryRun, true
}

// HasDryRun returns a boolean if a field has been set.
func (o *DeleteListenerRuleRequest) HasDryRun() bool {
	if o != nil && o.DryRun != nil {
		return true
	}

	return false
}

// SetDryRun gets a reference to the given bool and assigns it to the DryRun field.
func (o *DeleteListenerRuleRequest) SetDryRun(v bool) {
	o.DryRun = &v
}

// GetListenerRuleName returns the ListenerRuleName field value
func (o *DeleteListenerRuleRequest) GetListenerRuleName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ListenerRuleName
}

// GetListenerRuleNameOk returns a tuple with the ListenerRuleName field value
// and a boolean to check if the value has been set.
func (o *DeleteListenerRuleRequest) GetListenerRuleNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ListenerRuleName, true
}

// SetListenerRuleName sets field value
func (o *DeleteListenerRuleRequest) SetListenerRuleName(v string) {
	o.ListenerRuleName = v
}

func (o DeleteListenerRuleRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DryRun != nil {
		toSerialize["DryRun"] = o.DryRun
	}
	if true {
		toSerialize["ListenerRuleName"] = o.ListenerRuleName
	}
	return json.Marshal(toSerialize)
}

type NullableDeleteListenerRuleRequest struct {
	value *DeleteListenerRuleRequest
	isSet bool
}

func (v NullableDeleteListenerRuleRequest) Get() *DeleteListenerRuleRequest {
	return v.value
}

func (v *NullableDeleteListenerRuleRequest) Set(val *DeleteListenerRuleRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableDeleteListenerRuleRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableDeleteListenerRuleRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDeleteListenerRuleRequest(val *DeleteListenerRuleRequest) *NullableDeleteListenerRuleRequest {
	return &NullableDeleteListenerRuleRequest{value: val, isSet: true}
}

func (v NullableDeleteListenerRuleRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDeleteListenerRuleRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
