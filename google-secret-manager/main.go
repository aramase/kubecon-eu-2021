package main

import (
	"context"
	"flag"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"k8s.io/klog/v2"
)

var (
	secretName = flag.String("secret-name", "projects/csi-secret-demo-kubecon2021/secrets/app-secret/versions/latest", "secret name")
)

func main() {
	ctx := context.Background()

	// secret name is required
	if len(*secretName) == 0 {
		klog.Fatalf("secret name is required")
	}

	// create the client with the default credentials as token source
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		klog.Fatalf("failed to setup client: %+v", err)
	}

	// Build the request to access the secret
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: *secretName,
	}

	// call the API
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		klog.Fatalf("failed to access secret version: %v", err)
	}
	// log the secret
	klog.Info(string(result.Payload.Data))
	time.Sleep(3600 * time.Second)
}
