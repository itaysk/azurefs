package azfuse

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/itaysk/azurefs/azfuse"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

type AzureClientMock struct {
}

func NewAzureClientMock() azfuse.IAzureClient {
	log.Debugf("New Azure client mock")
	res := AzureClientMock{}
	return res
}

func (this AzureClientMock) GetAllResourceGroups() (*[]resources.Group, error) {
	res := []resources.Group{}
	res = append(res, resources.Group{Name: to.StringPtr("rg1")})
	res = append(res, resources.Group{Name: to.StringPtr("rg2")})
	return &res, nil
}

func (this AzureClientMock) GetAllResourcesInGroup(rgName string) (*[]resources.GenericResource, error) {
	res := []resources.GenericResource{}
	res = append(res, resources.GenericResource{Name: to.StringPtr("r1"), ID: to.StringPtr("r1")})
	res = append(res, resources.GenericResource{Name: to.StringPtr("r2"), ID: to.StringPtr("r2")})
	return &res, nil
}

func (this AzureClientMock) GetResourceJson(id string) ([]byte, error) {
	json := fmt.Sprintf(`{
id: "%s",
name: "%s"
}`, id, id)
	return []byte(json), nil
}

func (this AzureClientMock) GetTags() (*[]resources.TagDetails, error) {
	res := []resources.TagDetails{}
	res = append(res, resources.TagDetails{TagName: to.StringPtr("t1")})
	res = append(res, resources.TagDetails{TagName: to.StringPtr("t2")})
	return &res, nil
}

func (this AzureClientMock) FindAllByTag(tag string) (*[]resources.Group, *[]resources.GenericResource, error) {
	res1 := []resources.Group{}
	res1 = append(res1, resources.Group{Name: to.StringPtr("rg1")})
	res2 := []resources.GenericResource{}
	res2 = append(res2, resources.GenericResource{Name: to.StringPtr("r2"), ID: to.StringPtr("r2")})
	return &res1, &res2, nil
}
