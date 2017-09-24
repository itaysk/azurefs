package azfuse

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/itaysk/azurefs/azfuse"
	//"github.com/itaysk/azurefs/azfuse"
	log "github.com/sirupsen/logrus"
)

func TestAzureFs(t *testing.T) {
	//setup
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	mountPoint, err := ioutil.TempDir("", "TestMounting")
	if err != nil {
		t.Fatalf("couldn't create a temp directory : %s", err)
	}
	defer os.RemoveAll(mountPoint)

	azureClient := NewAzureClientMock()
	root := &azfuse.SubscriptionNode{Node: nodefs.NewDefaultNode()}
	fs := azfuse.NewAzureFs(azureClient, root)

	//test mount
	server, _, err := nodefs.MountRoot(mountPoint, fs.Root(), nil)
	defer server.Unmount()
	if err != nil {
		t.Fatalf("Mount fail: %v\n", err)
	}
	go server.Serve()

	if err := server.WaitMount(); err != nil {
		t.Fatal("WaitMount", err)
	}
	defer server.Unmount()

	//test rg
	rgs, err := ioutil.ReadDir(mountPoint)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	t.Log("resource groups:")
	for _, rg := range rgs {
		t.Log(rg.Name())
		if !rg.IsDir() {
			t.Logf("rg %s is not a directory", rg.Name())
			t.Fail()
		}
	}

	//test r
	rgPath := fmt.Sprintf("%s/%s", mountPoint, rgs[1].Name())
	rs, err := ioutil.ReadDir(rgPath)
	if err != nil {
		t.Logf("ReadDir error: %v", err)
		t.Fail()
	}
	t.Log("resources:")
	for _, r := range rs {
		t.Log(r.Name())
	}
	rPath := fmt.Sprintf("%s/%s", rgPath, rs[0].Name())
	f, err := ioutil.ReadFile(rPath)
	if err != nil {
		t.Logf("ReadFile error: %v", err)
		t.Fail()
	}
	t.Log(string(f))
}
