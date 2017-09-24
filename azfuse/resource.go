package azfuse

import (
	"bytes"
	"encoding/json"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	log "github.com/sirupsen/logrus"
)

type ResourceNode struct {
	nodefs.Node
	fs       *AzureFs
	Name     string
	Id       string
	size     uint64
	contents []byte
}

func (this *ResourceNode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {
	log.Debugf("in GetAttr for %s", this.Name)
	//TODO: not sure if can skip netowrk call here and set size to max (prints multiple times)
	err := this.cacheFile()
	if err != nil {
		log.Errorf("failed getting resource json for %s", this.Id)
		return fuse.OK //TODO: should return OK?
	}
	out.Size = this.size
	out.Mode = fuse.S_IFREG | 0444
	return fuse.OK
}

func (this *ResourceNode) Open(flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Debugf("in Open for %s", this.Name)
	err := this.cacheFile()
	if err != nil {
		log.Errorf("failed getting resource json for %s", this.Id)
		return nil, fuse.OK //TODO: should return OK?
	}
	return nodefs.NewDataFile(this.contents), fuse.OK
}

// cacheFile generates a JSON representation for the JSON from Azure API, and caches it in memory.
// it's safe to call before each access to file's properties
func (this *ResourceNode) cacheFile() error {
	if this.contents == nil {
		t, err := this.fs.azureClient.GetResourceJson(this.Id)
		if err != nil { //error should be handled by caller
			return err
		}
		var b bytes.Buffer
		json.Indent(&b, t, "", "    ")
		b.WriteByte('\n')
		this.size = uint64(b.Len())
		this.contents = b.Bytes()
	}
	return nil
}
