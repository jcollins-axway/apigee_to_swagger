package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-openapi/spec"
)

func main() {
	proxyDir := "/home/ubuntu/go/src/github.com/jcollins-axway/apigee_to_swagger/apiproxy/"
	proxyName := "Petstore"
	// Open our proxyFile
	proxyFile, err := os.Open(proxyDir + "/" + proxyName + ".xml")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(proxyFile)
	proxyFile.Close()

	// Unmarshal the proxy details
	var proxy APIProxy
	xml.Unmarshal(byteValue, &proxy)

	// Create teh paths object

	swagger := spec.SwaggerProps{
		BasePath: proxy.Basepaths,
		Info: &spec.Info{
			InfoProps: spec.InfoProps{
				Description: proxy.Description,
				Version:     proxy.Version.Major + "." + proxy.Version.Minor,
			},
		},
		Paths: &spec.Paths{
			Paths: map[string]spec.PathItem{},
		},
	}

	paths := swagger.Paths

	// Load the endpoint files
	for _, endpointFilename := range proxy.ProxyEndpoints.ProxyEndpoint {
		endpointFile, err := os.Open(proxyDir + "/proxies/" + endpointFilename + ".xml")
		if err != nil {
			fmt.Println(err)
		}

		byteValue, _ := ioutil.ReadAll(endpointFile)
		endpointFile.Close()

		// Unmarshal the proxy details
		var endpoint ProxyEndpoint
		xml.Unmarshal(byteValue, &endpoint)

		for _, flow := range endpoint.Flows.Flow {
			var verb, urlPath string
			operation := spec.Operation{
				OperationProps: spec.OperationProps{
					ID:          flow.Name,
					Description: flow.Description,
					Summary:     flow.Description,
				},
			}
			for _, condition := range flow.Conditions.Condition {
				if condition.Variable == "proxy.pathsuffix" && condition.Operator == "MatchesPath" {
					urlPath = condition.Value
					// Split path
					pathComponents := strings.Split(urlPath, "/")
					for i, pathComponent := range pathComponents {
						if pathComponent == "*" {
							// This is a * part of the url, change it to a variable name based on previous component
							pathComponents[i] = "{" + pathComponents[i-1] + "Id}"
						}
					}
					urlPath = strings.Join(pathComponents, "/")
				} else if condition.Variable == "request.verb" && (condition.Operator == "=" || condition.Operator == "equal") {
					verb = condition.Value

				}
			}
			pathProps := spec.PathItemProps{}
			if currentProps, ok := paths.Paths[urlPath]; ok {
				pathProps = currentProps.PathItemProps
			}

			switch strings.ToUpper(verb) {
			case "GET":
				pathProps.Get = &operation
			case "PUT":
				pathProps.Put = &operation
			case "POST":
				pathProps.Post = &operation
			case "DELETE":
				pathProps.Delete = &operation
			}

			paths.Paths[urlPath] = spec.PathItem{
				PathItemProps: pathProps,
			}

		}
		proxy.parsedEndpoints = append(proxy.parsedEndpoints, endpoint)
	}

	swaggerBytes, _ := json.Marshal(swagger)

	fmt.Println(string(swaggerBytes))
}
