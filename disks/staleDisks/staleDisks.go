package staleDisks

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func GetUnusedDisks(ctx context.Context, projectID string) []string {
	// Get a disk client so we could perform related operations
	diskClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		panic(err)

	}
	defer diskClient.Close()

	// Protobuf that acts as a filter and other options that we may want to implement
	req := &computepb.AggregatedListDisksRequest{
		Project: projectID,
	}

	// Aggregated List gives an iterator that fetches data from all Regions so that
	// we don't have to loop over every region
	it := diskClient.AggregatedList(ctx, req)

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
			for _, instance := range disks {
				if len(instance.Users) == 0 {
					unusedDisks = append(unusedDisks, *instance.Name)
				}
			}
		}
	}

	return unusedDisks
}
