package azureHelper

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

// GetLatestAPIVersion looks up the latest API Version for the specified resource provider and resource type
func GetLatestAPIVersion(client *resources.ProvidersClient, provider string, resourceType string) (result string) {
	providerRes, _ := client.Get(provider, "resourceTypes/apiVersion")
	for _, p := range *providerRes.ResourceTypes {
		if *p.ResourceType == resourceType {
			result = (*p.APIVersions)[0]
		}
	}
	log.Printf("latest api version for provider %s, resourceType %s, is %s", provider, resourceType, result)
	return
}

// GetLatestAPIVersionByID looks up the latest API Version for the specified resource id
func GetLatestAPIVersionByID(client *resources.ProvidersClient, id string) (result string) {
	idSplit := strings.Split(id, "/")
	result = GetLatestAPIVersion(client, idSplit[6], idSplit[7])
	return
}
