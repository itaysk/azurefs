package azfuse

import (
	log "github.com/sirupsen/logrus"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type TagNode struct {
	nodefs.Node
	Name string
}

func (this *TagNode) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	log.Debugf("in OpenDir for tag %s", this.Name)
	res := []fuse.DirEntry{}
	rgs, rs, err := azureClient.FindAllByTag(this.Name)
	if err != nil {
		log.Fatalf("Failed finding tags for %s", this.Name)
		log.Fatal(err)
		return res, fuse.OK
	}

	for _, r := range *rs {
		res = append(res, fuse.DirEntry{Name: *r.Name})
		rn := ResourceNode{Node: nodefs.NewDefaultNode(), Name: *r.Name, Id: *r.ID}
		this.Inode().NewChild(rn.Name, false, &rn)
	}
	for _, rg := range *rgs {
		res = append(res, fuse.DirEntry{Name: *rg.Name})
		rgn := ResourceGroupNode{Node: nodefs.NewDefaultNode(), Name: *rg.Name}
		this.Inode().NewChild(*rg.Name, true, &rgn)

		rs, err := azureClient.GetAllResourcesInGroup(*rg.Name)
		if err != nil {
			log.Fatalf("Failed getting resources for group %s", *rg.Name)
			log.Fatal(err)
			return res, fuse.OK
		}
		for _, r := range *rs {
			rn := ResourceNode{Node: nodefs.NewDefaultNode(), Name: *r.Name, Id: *r.ID}
			rgn.Inode().NewChild(rn.Name, false, &rn)
		}
	}
	return res, fuse.OK
}
