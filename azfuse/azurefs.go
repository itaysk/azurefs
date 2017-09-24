package azfuse

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type AzureFs struct {
	azureClient IAzureClient
	root        *SubscriptionNode
}

func NewAzureFs(azureClient IAzureClient, root nodefs.Node) *AzureFs {

	azfs := &AzureFs{
		root:        root.(*SubscriptionNode),
		azureClient: azureClient,
	}
	azfs.root.fs = azfs
	return azfs
}

func (fs *AzureFs) String() string {
	return "azurefs"
}

func (fs *AzureFs) Root() nodefs.Node {
	return fs.root
}
