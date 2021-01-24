package main

import (
	"encoding/xml"
	"regexp"
	"strings"

	"github.com/oriser/regroup"
)

// APIProxy structures

// APIProxy - APIGEE API Proxy definition
type APIProxy struct {
	Basepaths       string               `xml:"Basepaths"`
	Version         ConfigurationVersion `xml:"ConfigurationVersion"`
	CreatedAt       int                  `xml:"CreatedAt"`
	CreatedBy       string               `xml:"CreatedBy"`
	Description     string               `xml:"Description"`
	DisplayName     string               `xml:"DisplayName"`
	LastModifiedAt  int                  `xml:"LastModifiedAt"`
	LastModifiedBy  string               `xml:"LastModifiedBy"`
	ManifestVersion string               `xml:"ManifestVersion"`
	Policies        Policies             `xml:"Policies"`
	ProxyEndpoints  ProxyEndpoints       `xml:"ProxyEndpoints"`
	Resources       string               `xml:"Resources"`
	Spec            string               `xml:"Spec"`
	TargetServers   string               `xml:"TargetServers"`
	TargetEndpoints string               `xml:"TargetEndpoints"`
	parsedEndpoints []ProxyEndpoint
}

// ConfigurationVersion - APIGEE API Proxy version
type ConfigurationVersion struct {
	Major string `xml:"majorVersion,attr"`
	Minor string `xml:"minorVersion,attr"`
}

// Policies - APIGEE API Proxy Policy filenames
type Policies struct {
	Policies []string `xml:"Policy"`
}

// ProxyEndpoints - APIGEE API Proxy Endpoint filenames
type ProxyEndpoints struct {
	ProxyEndpoint []string `xml:"ProxyEndpoint"`
}

// ProxyEndpoint - APIGEE Proxy Endpoint file structure
type ProxyEndpoint struct {
	Name    string  `xml:"name,attr"`
	PreFlow PreFlow `xml:"PreFlow"`
	Flows   Flows   `xml:"Flows"`
}

// PreFlow - APIGEE Proxy endpoint preflow
type PreFlow struct {
	Name     string   `xml:"name,attr"`
	Request  Request  `xml:"Request"`
	Response Response `xml:"Response"`
}

// Request - APIGEE Proxy request steps
type Request struct {
	Step []Step `xml:"Step"`
}

// Response - APIGEE Proxy response steps
type Response struct {
	Step []Step `xml:"Step"`
}

// Step - APIGEE Proxy step names
type Step struct {
	Name string `xml:"Name"`
}

// Flows - APIGEE Proxy flows
type Flows struct {
	Flow []Flow `xml:"Flow"`
}

// Flow - APIGEE proxy flow
type Flow struct {
	Name        string     `xml:"name,attr"`
	Description string     `xml:"Description"`
	Request     Request    `xml:"Request"`
	Response    Response   `xml:"Response"`
	Conditions  Conditions `xml:"Condition"`
}

// Conditions - array of conditions
type Conditions struct {
	Condition []Condition
}

// Condition -- parsed representation of APIGEE condition
type Condition struct {
	Variable string
	Operator string
	Value    string
}

//UnmarshalXML - custom XML unmarshall to parse APIGEE flow conditions
func (c *Conditions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var el string

	if err := d.DecodeElement(&el, &start); err != nil {
		return err
	}

	if el == "null" {
		return nil
	}

	// Split all conditions
	re := regexp.MustCompile("\\).*\\(")
	conditions := re.Split(el, -1)

	// iterate conditions
	for _, condition := range conditions {
		r := regroup.MustCompile("^\\({0,1}(?P<var>[a-z\\.]*) (?P<op>.*) (?P<val>.*\")\\){0,1}$")

		matches, err := r.Groups(condition)
		if err != nil {
			return err
		}

		// add each condition to the array
		c.Condition = append(c.Condition, Condition{
			Variable: strings.Trim(matches["var"], "\""),
			Operator: strings.Trim(matches["op"], "\""),
			Value:    strings.Trim(matches["val"], "\""),
		})
	}

	return nil
}
