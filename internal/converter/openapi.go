package converter

import (
	"encoding/json"
	"github.com/alecthomas/jsonschema"
)

// OpenAPI is the root document object of the OpenAPI document
type OpenAPI struct {
	VendorExtensible `json:"-"`

	Components   Components      `json:"components,omitempty"`
	ExternalDocs json.RawMessage `json:"externalDocs,omitempty"`
	Info         json.RawMessage `json:"info,omitempty"`
	OpenAPI      string          `json:"openapi"`
	Paths        json.RawMessage `json:"paths,omitempty"`
	Security     json.RawMessage `json:"security,omitempty"`
	Servers      json.RawMessage `json:"servers,omitempty"`
	Tags         json.RawMessage `json:"tags,omitempty"`
}

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
type Components struct {
	VendorExtensible `json:"-"`

	Schemas         Schemas         `json:"schemas,omitempty"`
	Responses       json.RawMessage `json:"responses,omitempty"`
	Parameters      json.RawMessage `json:"parameters,omitempty"`
	Examples        json.RawMessage `json:"examples,omitempty"`
	RequestBodies   json.RawMessage `json:"requestBodies,omitempty"`
	Headers         json.RawMessage `json:"headers,omitempty"`
	SecuritySchemes json.RawMessage `json:"securitySchemes,omitempty"`
	Links           json.RawMessage `json:"links,omitempty"`
	Callbacks       json.RawMessage `json:"callbacks,omitempty"`
}

type Schemas map[string]*jsonschema.Type
type VendorExtensible map[string]json.RawMessage
