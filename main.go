package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/namsral/flag"

	"github.com/itaysk/azurefs/azfuse"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

func main() {
	//flag settings (note we're using github.com/namsral/flag instead of standard flag which automatically merges with environment variables)
	var mountPoint = flag.String("mount-point", "", "an empty directory to mount at")
	var azureSubscriptionId = flag.String("azure-subscription-id", "", "subscription id to mount. Alternatively, set environment variable 'AZURE_SUBSCRIPTION_ID'")
	var azureClientId = flag.String("azure-client-id", "", "app id of a spn with required permissions. Alternatively, set environment variable 'AZURE_CLIENT_ID'")
	var azureClientSecret = flag.String("azure-client-secret", "", "app key of a spn with required permissions. Alternatively, set environment variable 'AZURE_CLIENT_SECRET'")
	var azureTenantId = flag.String("azure-tenant-id", "", "Id of the Azure AD that is associated with the subscription. Alternatively, set environment variable 'AZURE_TENANT_ID'")
	var isVerbose = flag.Bool("v", false, "Print verbose (debug) level messages")
	flag.Parse()

	if *isVerbose {
		log.SetLevel(log.DebugLevel)
	}

	azureSettings := map[string]string{
		"subscriptionId": *azureSubscriptionId,
		"clientId":       *azureClientId,
		"clientSecret":   *azureClientSecret,
		"tenantId":       *azureTenantId,
	}

	//input validation
	isInputValid := true
	//are all azure settings args provided?
	if azureSettings["subscriptionId"] == "" ||
		azureSettings["clientId"] == "" ||
		azureSettings["clientSecret"] == "" ||
		azureSettings["tenantId"] == "" {
		isInputValid = false
		log.Errorf("missing azure setting. current settings: %s", azureSettings)
	}
	//is mount point arg provided?
	if *mountPoint == "" {
		isInputValid = false
		log.Errorln("missing mount point")

	} else {
		//is the mount point provided valid?
		if _, err := os.Stat(*mountPoint); os.IsNotExist(err) {
			isInputValid = false
			log.Errorf("provided mount point: %s is invalid", *mountPoint)

		}
	}

	//if validation failed, print usage details
	if !isInputValid {
		fmt.Printf("usage: %s --mount-point dir \n", os.Args[0])
		fmt.Println("Once the FUSE server is running it will block, so consider running this in background")
		flag.PrintDefaults()
		os.Exit(2)
	}

	//create and start the fuse server
	azureClient := azfuse.NewAzureClient(azureSettings)
	root := &azfuse.SubscriptionNode{Node: nodefs.NewDefaultNode()}
	fs := azfuse.NewAzureFs(azureClient, root)
	log.Infof("mounting on %s", *mountPoint)
	server, _, err := nodefs.MountRoot(*mountPoint, fs.Root(), nil)
	defer server.Unmount()
	if err != nil {
		log.Errorf("Mount fail: %v\n", err)
		os.Exit(1)
	}
	handleSigint(server, *mountPoint)
	server.Serve()
	//todo: waitmount?

	//main is expected to never return, since server.Serve() will block.
}

// used to capture sigint (ctrl+c), and gracefully unmount
func handleSigint(srv *fuse.Server, mountpoint string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		log.Info("Unmounting")
		err := srv.Unmount()
		if err != nil {
			log.Error(err)
			log.Debug("Trying lazy unmount")
			cmd := exec.Command("fusermount", "-u", "-z", mountpoint) //consider falling back to sudo umount mount-point
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
		os.Exit(1)
	}()
}
