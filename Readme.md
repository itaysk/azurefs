Mount your Azure account as a filesystem, so you could `cd` into your Azure account, `ls` all resource groups, `cd` into a resource group, and `cat` a resource.  
Reuse your Linux productivity skills and tools to script against Azure like regular files and directories.

## Structure

- Subscription (root)
      - Resource Groups
            - Resources (json)
      - @tags
            - Tags
                  - Resources (json)
                  - Resource Groups
                        - Resources (json)

All nodes are directories, except leaf nodes which are files
First level directories that start with "!" are artificial directories injected to provide additional functionality (Resource Group name can't start with "@" so you can be sure those won't conflict)


## Known issues
- Network calls on GetAttr cause `ls` in Resource Group level to get and cache every resource in that resource group ([more info[(https://stackoverflow.com/questions/46267972/fuse-avoid-calculating-size-in-getattr)]).
- Only basic FUSE callbacks are implemented.
- You have to open a tag folder in order to populate it's content. You can't `cd @tags/someTag/someRg` before you `cd` into 'someTag' because 'someRg' won't exist yet in this path.
- Only tag names are considered, not tag values.
- All Azure operations are limited to return 50 items

## Usage

```
usage: ./azfuse --mount-point dir
Once the FUSE server is running it will block, so consider running this in background
  -azure-client-id string
        app id of a spn with required permissions. Alternatively, set environment variable 'AZURE_CLIENT_ID'
  -azure-client-secret string
        app key of a spn with required permissions. Alternatively, set environment variable 'AZURE_CLIENT_SECRET'
  -azure-subscription-id string
        subscription id to mount. Alternatively, set environment variable 'AZURE_SUBSCRIPTION_ID'
  -azure-tenant-id string
        Id of the Azure AD that is associated with the subscription. Alternatively, set environment variable 'AZURE_TENANT_ID'
  -mount-point string
        an empty directory to mount at
  -v    Print verbose (debug) level messages
```

## Built on

- Ubuntu 16.04.3 LTS
- Go 1.6.2
- https://github.com/hanwen/go-fuse (master)

## Contribution

Please do! Open an issue describing the fix or feature, once validated feel free to submit a PR.