package azfuse

import (
	"encoding/json"
	"fmt"

	"github.com/itaysk/azurefs/azureHelper"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

type AzureClient struct {
	autorestClient  autorest.Client
	groupsClient    resources.GroupsClient
	groupClient     resources.GroupClient
	resourcesClient azureHelper.ResourcesClient
	providersClient resources.ProvidersClient
	tagsClient      resources.TagsClient
}

var subscriptionID string
var pageSize int32 = 50

func NewAzureClient(azureSettings map[string]string) AzureClient {
	log.Debugf("New Azure client with settings: +%v", azureSettings)
	subscriptionID = azureSettings["subscriptionId"]
	res := AzureClient{}
	auth, err := getAuthorizer(azureSettings)
	if err != nil {
		log.Error("Failed creating autorest authorizer")
		log.Fatal(err)
	}
	autorestClient := &autorest.Client{}
	autorestClient.Authorizer = auth
	res.groupsClient = resources.NewGroupsClient(subscriptionID)
	res.groupsClient.Authorizer = auth
	res.groupClient = resources.NewGroupClient(subscriptionID)
	res.groupClient.Authorizer = auth
	res.providersClient = resources.NewProvidersClient(subscriptionID)
	res.providersClient.Authorizer = auth
	res.resourcesClient = azureHelper.ResourcesClient{&res.groupClient}
	res.tagsClient = resources.NewTagsClient(subscriptionID)
	res.tagsClient.Authorizer = auth
	return res
}

// based on Azure-Autorest's `util.GetAuthorizer`: https://github.com/Azure/go-autorest/blob/master/autorest/utils/auth.go
// modified to not rely only on hardcoded env vars
func getAuthorizer(azureSettings map[string]string) (*autorest.BearerAuthorizer, error) {
	aadEndpoint := "https://login.microsoftonline.com/"
	armEndpoint := "https://management.core.windows.net/"

	oauthConfig, err := adal.NewOAuthConfig(aadEndpoint, azureSettings["tenantId"])
	if err != nil { //error should be handled by caller
		return nil, err
	}

	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, azureSettings["clientId"], azureSettings["clientSecret"], armEndpoint)
	if err != nil { //error should be handled by caller
		return nil, err
	}

	return autorest.NewBearerAuthorizer(spToken), nil
}

func (this AzureClient) GetAllResourceGroups() (*[]resources.Group, error) {
	return this.getResourceGroups("")
}

func (this AzureClient) getResourceGroups(filter string) (*[]resources.Group, error) {
	log.Debugf("getting resource groups, filter %s", filter)
	groups, err := this.groupsClient.List(filter, &pageSize)
	if err != nil { //error should be handled by caller
		return nil, err
	}
	log.Debugf("got %d groups", len(*groups.Value))

	return groups.Value, nil
}

func (this AzureClient) GetAllResourcesInGroup(rgName string) (*[]resources.GenericResource, error) {
	log.Debugf("getting resources for rg %s", rgName)
	resources, err := this.groupsClient.ListResources(rgName, "", "", &pageSize)
	if err != nil { //error should be handled by caller
		return nil, err
	}
	log.Debugf("got %d resources", len(*resources.Value))

	return resources.Value, nil
}

func (this AzureClient) GetResourceJson(id string) ([]byte, error) {
	log.Debugf("getting json for %s", id)
	apiVersion := azureHelper.GetLatestAPIVersionByID(&this.providersClient, id)
	r, err := this.resourcesClient.GetByID(id, apiVersion)
	if err != nil { //error should be handled by caller
		return nil, err
	}
	rJson, err := json.Marshal(r)
	if err != nil {
		log.Error("Failed to serialize resource to JSON")
		return nil, err
	}
	log.Debugf("got %d bytes", len(rJson))
	return rJson, nil
}

func (this AzureClient) GetTags() (*[]resources.TagDetails, error) {
	log.Debug("getting tags")
	tags, err := this.tagsClient.List()
	if err != nil { //error should be handled by caller
		return nil, err
	}
	log.Debugf("got %d tags", len(*tags.Value))

	return tags.Value, nil
}

func (this AzureClient) FindAllByTag(tag string) (*[]resources.Group, *[]resources.GenericResource, error) {
	log.Debugf("finding all for tag %s", tag)
	filter := fmt.Sprintf("tagname eq '%s'", tag)
	rgs, err := this.getResourceGroups(filter)
	if err != nil { //error should be handled by caller
		return nil, nil, err
	}
	rs, err := this.getResources(filter)
	if err != nil { //error should be handled by caller
		return nil, nil, err
	}
	log.Debugf("found %d", len(*rgs)+len(*rs))
	//TODO: consider generalizing into in interface, and returning a single value?
	return rgs, rs, nil
}

func (this AzureClient) getResources(filter string) (*[]resources.GenericResource, error) {
	log.Debugf("getting resources, filter %s", filter)
	rs, err := this.groupClient.List(filter, "", &pageSize)
	if err != nil { //error should be handled by caller
		return nil, err
	}
	log.Debugf("got %d resources", len(*rs.Value))

	return rs.Value, nil
}
