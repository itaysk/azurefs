package azfuse

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type AzureFs struct {
	root *SubscriptionNode
}

var azureClient AzureClient

func NewAzureFs(azureSettings map[string]string) *AzureFs {
	azureClient = NewAzureClient(azureSettings)

	return &AzureFs{
		root: &SubscriptionNode{Node: nodefs.NewDefaultNode()},
	}
}

func (fs *AzureFs) String() string {
	return "azurefs"
}

func (fs *AzureFs) Root() nodefs.Node {
	return fs.root
}
