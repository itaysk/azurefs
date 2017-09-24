package azfuse

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type ResourceGroupNode struct {
	nodefs.Node
	fs   *AzureFs
	Name string
}
