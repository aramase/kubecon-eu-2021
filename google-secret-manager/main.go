package main

import (
	"context"
	"flag"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"k8s.io/klog/v2"
)

var (
	secretName = flag.String("secret-name", "app-secret", "secret name")
)

func main() {
	ctx := context.Background()

	// secret name is required
	if len(*secretName) == 0 {
		klog.Fatalf("secret name is required")
	}

	// get the credentials to access secret manager
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		klog.Fatalf("failed to get credentials: %v", err)
	}

	// create the client with the default credentials as token source
	client, err := secretmanager.NewClient(ctx, []option.ClientOption{option.WithTokenSource(creds.TokenSource)}...)
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
}
