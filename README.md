# Client
v0.1

## client is a package for test
### example
##### Create workload and Client
```go
// Create client
c := client.NewClient()

// Get yaml files for workload
f := client.GetTestFileListToBytes(client.TestFilePath, client.Pvc, client.Deployment, client.Service)

// Create workload object
newWorkload := client.NewWorkload("cdm-test", f)

// Add workload to client
c.AddWorkload(newWorkload)
```
##### Apply workload
```go
// Apply workload
c.ApplyWorkload("cdm-test")
```
##### Delete workload
```go
// Delete workload
c.DeleteWorkload("cdm-test")
```