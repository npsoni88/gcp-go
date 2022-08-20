package staleDisks

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func GetUnusedDisks(ctx context.Context, projectID string) []string {
	//instancesClient, err := compute.NewInstancesRESTClient(ctx)
	diskClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		panic(err)

	}
	defer diskClient.Close()

	// Use the `MaxResults` parameter to limit the number of results that the API returns per response page.
	req := &computepb.AggregatedListDisksRequest{
		Project: projectID,
	}

	it := diskClient.AggregatedList(ctx, req)
	fmt.Printf("Instances found:\n")

	var unusedDisks []string
	for {
		disk, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		disks := disk.Value.Disks
		if len(disks) > 0 {
			//fmt.Printf("zone -> %s\n", disk.Key)
			for _, instance := range disks {
				//fmt.Printf("disk - %s with size - %dGb has %d users\n", *instance.Name, *instance.SizeGb, len(instance.Users))
				if len(instance.Users) == 0 {
					unusedDisks = append(unusedDisks, *instance.Name)
				}
			}
		}
	}

	return unusedDisks
}
