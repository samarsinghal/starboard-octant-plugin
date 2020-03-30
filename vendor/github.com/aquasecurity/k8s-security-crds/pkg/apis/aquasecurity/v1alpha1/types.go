package v1alpha1

import (
	"k8s.io/utils/pointer"
	"strconv"

	"github.com/aquasecurity/k8s-security-crds/pkg/apis/aquasecurity"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	VulnerabilitiesCRName    = "vulnerabilities.aquasecurity.github.com"
	VulnerabilitiesCRVersion = "v1alpha1"
)

var (
	VulnerabilitiesCRD = extv1beta1.CustomResourceDefinition{
		ObjectMeta: meta.ObjectMeta{
			Name: VulnerabilitiesCRName,
		},
		Spec: extv1beta1.CustomResourceDefinitionSpec{
			Group: aquasecurity.GroupName,
			Versions: []extv1beta1.CustomResourceDefinitionVersion{
				{
					Name:    VulnerabilitiesCRVersion,
					Served:  true,
					Storage: true,
				},
			},
			Scope: extv1beta1.NamespaceScoped,
			Names: extv1beta1.CustomResourceDefinitionNames{
				Singular:   "vulnerability",
				Plural:     "vulnerabilities",
				Kind:       "Vulnerability",
				ListKind:   "VulnerabilityList",
				Categories: []string{"all"},
				ShortNames: []string{"vulns", "vuln"},
			},
			Validation: &extv1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &extv1beta1.JSONSchemaProps{
					Type: "object",
					Required: []string{
						"apiVersion",
						"kind",
						"metadata",
						"report",
					},
					Properties: map[string]extv1beta1.JSONSchemaProps{
						"apiVersion": {Type: "string"},
						"kind":       {Type: "string"},
						"metadata":   {Type: "object"},
						"report": {
							Type: "object",
							Required: []string{
								"generatedAt",
								"scanner",
								"artifact",
								"vulnerabilities",
							},
							Properties: map[string]extv1beta1.JSONSchemaProps{
								"generatedAt": {
									Type:   "string",
									Format: "date-time",
								},
								"scanner": {
									Type: "object",
									Required: []string{
										"name",
										"vendor",
										"version",
									},
									Properties: map[string]extv1beta1.JSONSchemaProps{
										"name":    {Type: "string"},
										"vendor":  {Type: "string"},
										"version": {Type: "string"},
									},
								},
								"registry": {
									Type: "object",
									Properties: map[string]extv1beta1.JSONSchemaProps{
										"url": {Type: "string", Format: "url"},
									},
								},
								"artifact": {
									Type: "object",
									Properties: map[string]extv1beta1.JSONSchemaProps{
										"repository": {Type: "string"},
										"digest":     {Type: "string"},
										"tag":        {Type: "string"},
										"mimeType":   {Type: "string"},
									},
								},
								"summary": {
									Type: "object",
									Required: []string{
										"criticalCount",
										"highCount",
										"mediumCount",
										"lowCount",
										"unknownCount",
									},
									Properties: map[string]extv1beta1.JSONSchemaProps{
										"criticalCount": {Type: "integer", Minimum: pointer.Float64Ptr(0)},
										"highCount":     {Type: "integer", Minimum: pointer.Float64Ptr(0)},
										"mediumCount":   {Type: "integer", Minimum: pointer.Float64Ptr(0)},
										"lowCount":      {Type: "integer", Minimum: pointer.Float64Ptr(0)},
										"unknownCount":  {Type: "integer", Minimum: pointer.Float64Ptr(0)},
									},
								},
								"vulnerabilities": {
									Type: "array",
									Items: &extv1beta1.JSONSchemaPropsOrArray{
										Schema: &extv1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"vulnerabilityID",
												"resource",
												"installedVersion",
												"fixedVersion",
												"severity",
												"title",
											},
											Properties: map[string]extv1beta1.JSONSchemaProps{
												"vulnerabilityID":  {Type: "string"},
												"resource":         {Type: "string"},
												"installedVersion": {Type: "string"},
												"fixedVersion":     {Type: "string"},
												"severity": {
													Type: "string",
													Enum: []extv1beta1.JSON{
														{Raw: []byte(strconv.Quote(string(SeverityCritical)))},
														{Raw: []byte(strconv.Quote(string(SeverityHigh)))},
														{Raw: []byte(strconv.Quote(string(SeverityMedium)))},
														{Raw: []byte(strconv.Quote(string(SeverityLow)))},
														{Raw: []byte(strconv.Quote(string(SeverityUnknown)))},
													},
												},
												"title":       {Type: "string"},
												"layerID":     {Type: "string"},
												"description": {Type: "string"},
												"links": {
													Type: "array",
													Items: &extv1beta1.JSONSchemaPropsOrArray{
														Schema: &extv1beta1.JSONSchemaProps{
															Type: "string",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
)

type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityUnknown  Severity = "UNKNOWN"
)

// VulnerabilityScanner is the spec for a vulnerability scanner.
type VulnerabilityScanner struct {
	Name    string `json:"name"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	// TODO Add generic or vendor specific properties such as com.github.aquasecurity/trivy-db-updated-at
}

type VulnerabilitySummary struct {
	CriticalCount int `json:"criticalCount"`
	HighCount     int `json:"highCount"`
	MediumCount   int `json:"mediumCount"`
	LowCount      int `json:"lowCount"`
	UnknownCount  int `json:"unknownCount"`
}

type Registry struct {
	URL string `json:"url"`
}

// Artifact is the spec for an artifact that can be scanned.
type Artifact struct {
	Repository string `json:"repository"`
	Digest     string `json:"digest"`
	Tag        string `json:"tag,omitempty"`
	MimeType   string `json:"mimeType,omitempty"`
}

// VulnerabilityItem is the spec for a vulnerability record.
type VulnerabilityItem struct {
	VulnerabilityID string `json:"vulnerabilityID"`
	Resource        string `json:"resource"`
	// TODO Add ResourceType enum property to distinguish between OS packages and application dependencies
	InstalledVersion string   `json:"installedVersion"`
	FixedVersion     string   `json:"fixedVersion"`
	Severity         Severity `json:"severity"`
	LayerID          string   `json:"layerID"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Links            []string `json:"links"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Vulnerability is a specification for the Vulnerability resource.
type Vulnerability struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Report VulnerabilityReport `json:"report"`
}

// VulnerabilityReport is the spec for the vulnerability report.
//
// The spec follows the Pluggable Scanners API defined for Harbor.
// @see https://github.com/goharbor/pluggable-scanner-spec/blob/master/api/spec/scanner-adapter-openapi-v1.0.yaml
type VulnerabilityReport struct {
	GeneratedAt     meta.Time            `json:"generatedAt"`
	Scanner         VulnerabilityScanner `json:"scanner"`
	Registry        Registry             `json:"registry"`
	Artifact        Artifact             `json:"artifact"`
	Summary         VulnerabilitySummary `json:"summary"`
	Vulnerabilities []VulnerabilityItem  `json:"vulnerabilities"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VulnerabilityList is a list of Vulnerability resources.
type VulnerabilityList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`

	Items []Vulnerability `json:"items"`
}
