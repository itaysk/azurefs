package azfuse

import (
	log "github.com/sirupsen/logrus"

	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type SubscriptionNode struct {
	nodefs.Node
}

func (this *SubscriptionNode) Deletable() bool {
	return false
}

func (this *SubscriptionNode) OnMount(c *nodefs.FileSystemConnector) {
	log.Debugf("mount called")

	rgs, err := azureClient.GetAllResourceGroups()
	if err != nil {
		log.Fatalf("Failed getting resource groups")
		return
	}
	for _, rg := range *rgs {
		rgn := ResourceGroupNode{Node: nodefs.NewDefaultNode(), Name: *rg.Name}
		this.Inode().NewChild(rgn.Name, true, &rgn)
		rs, err := azureClient.GetAllResourcesInGroup(*rg.Name)
		if err != nil {
			log.Error("Failed getting resource groups")
			return
		}
		for _, r := range *rs {
			rn := ResourceNode{Node: nodefs.NewDefaultNode(), Name: *r.Name, Id: *r.ID}
			rgn.Inode().NewChild(rn.Name, false, &rn)
		}
	}

	tagsContainer := nodefs.NewDefaultNode()
	this.Inode().NewChild("@tags", true, tagsContainer)
	tags, err := azureClient.GetTags()
	if err != nil {
		log.Error("Failed getting tags")
		return
	}
	for _, t := range *tags {
		tn := TagNode{Node: nodefs.NewDefaultNode(), Name: *t.TagName}
		tagsContainer.Inode().NewChild(tn.Name, true, &tn)
	}

}
