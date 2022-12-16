/*
caraml/timber/v1/dataset_service.proto

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: version not set
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// V1LogMetadata struct for V1LogMetadata
type V1LogMetadata struct {
	// Unique identifier of a log generated by a LogProducer.
	Id *string `json:"id,omitempty"`
	// Name of the log, generated by Dataset Service.
	Name *string `json:"name,omitempty"`
	Type *V1LogType `json:"type,omitempty"`
	// List of target names associated with a log.
	TargetNames []string `json:"targetNames,omitempty"`
	// BQ table ID where the data is stored.
	BqTable *string `json:"bqTable,omitempty"`
	LogProducer *V1LogProducer `json:"logProducer,omitempty"`
}

// NewV1LogMetadata instantiates a new V1LogMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1LogMetadata() *V1LogMetadata {
	this := V1LogMetadata{}
	var type_ V1LogType = V1LOGTYPE_UNSPECIFIED
	this.Type = &type_
	return &this
}

// NewV1LogMetadataWithDefaults instantiates a new V1LogMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1LogMetadataWithDefaults() *V1LogMetadata {
	this := V1LogMetadata{}
	var type_ V1LogType = V1LOGTYPE_UNSPECIFIED
	this.Type = &type_
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *V1LogMetadata) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1LogMetadata) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *V1LogMetadata) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *V1LogMetadata) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1LogMetadata) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1LogMetadata) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1LogMetadata) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1LogMetadata) SetName(v string) {
	o.Name = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *V1LogMetadata) GetType() V1LogType {
	if o == nil || o.Type == nil {
		var ret V1LogType
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1LogMetadata) GetTypeOk() (*V1LogType, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *V1LogMetadata) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given V1LogType and assigns it to the Type field.
func (o *V1LogMetadata) SetType(v V1LogType) {
	o.Type = &v
}

// GetTargetNames returns the TargetNames field value if set, zero value otherwise.
func (o *V1LogMetadata) GetTargetNames() []string {
	if o == nil || o.TargetNames == nil {
		var ret []string
		return ret
	}
	return o.TargetNames
}

// GetTargetNamesOk returns a tuple with the TargetNames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1LogMetadata) GetTargetNamesOk() ([]string, bool) {
	if o == nil || o.TargetNames == nil {
		return nil, false
	}
	return o.TargetNames, true
}

// HasTargetNames returns a boolean if a field has been set.
func (o *V1LogMetadata) HasTargetNames() bool {
	if o != nil && o.TargetNames != nil {
		return true
	}

	return false
}

// SetTargetNames gets a reference to the given []string and assigns it to the TargetNames field.
func (o *V1LogMetadata) SetTargetNames(v []string) {
	o.TargetNames = v
}

// GetBqTable returns the BqTable field value if set, zero value otherwise.
func (o *V1LogMetadata) GetBqTable() string {
	if o == nil || o.BqTable == nil {
		var ret string
		return ret
	}
	return *o.BqTable
}

// GetBqTableOk returns a tuple with the BqTable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1LogMetadata) GetBqTableOk() (*string, bool) {
	if o == nil || o.BqTable == nil {
		return nil, false
	}
	return o.BqTable, true
}

// HasBqTable returns a boolean if a field has been set.
func (o *V1LogMetadata) HasBqTable() bool {
	if o != nil && o.BqTable != nil {
		return true
	}

	return false
}

// SetBqTable gets a reference to the given string and assigns it to the BqTable field.
func (o *V1LogMetadata) SetBqTable(v string) {
	o.BqTable = &v
}

// GetLogProducer returns the LogProducer field value if set, zero value otherwise.
func (o *V1LogMetadata) GetLogProducer() V1LogProducer {
	if o == nil || o.LogProducer == nil {
		var ret V1LogProducer
		return ret
	}
	return *o.LogProducer
}

// GetLogProducerOk returns a tuple with the LogProducer field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1LogMetadata) GetLogProducerOk() (*V1LogProducer, bool) {
	if o == nil || o.LogProducer == nil {
		return nil, false
	}
	return o.LogProducer, true
}

// HasLogProducer returns a boolean if a field has been set.
func (o *V1LogMetadata) HasLogProducer() bool {
	if o != nil && o.LogProducer != nil {
		return true
	}

	return false
}

// SetLogProducer gets a reference to the given V1LogProducer and assigns it to the LogProducer field.
func (o *V1LogMetadata) SetLogProducer(v V1LogProducer) {
	o.LogProducer = &v
}

func (o V1LogMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.TargetNames != nil {
		toSerialize["targetNames"] = o.TargetNames
	}
	if o.BqTable != nil {
		toSerialize["bqTable"] = o.BqTable
	}
	if o.LogProducer != nil {
		toSerialize["logProducer"] = o.LogProducer
	}
	return json.Marshal(toSerialize)
}

type NullableV1LogMetadata struct {
	value *V1LogMetadata
	isSet bool
}

func (v NullableV1LogMetadata) Get() *V1LogMetadata {
	return v.value
}

func (v *NullableV1LogMetadata) Set(val *V1LogMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableV1LogMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableV1LogMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1LogMetadata(val *V1LogMetadata) *NullableV1LogMetadata {
	return &NullableV1LogMetadata{value: val, isSet: true}
}

func (v NullableV1LogMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1LogMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


