package main

import (
	"context"
	"flag"
	"fmt"
	"regexp"

	kv "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"k8s.io/klog/v2"
)

var (
	keyvaultName          = flag.String("keyvault-name", "", "Azure Keyvault name")
	keyvaultSecretName    = flag.String("keyvault-secret-name", "", "Azure Keyvault secret name")
	keyvaultSecretVersion = flag.String("keyvault-secret-version", "", "Azure Keyvault secret version")

	// credentials for accessing keyvault
	clientID     = flag.String("client-id", "", "Service principal clientID to access keyvault")
	clientSecret = flag.String("client-secret", "", "Service principal client secret to access keyvault")
	tenantID     = flag.String("tenant-id", "", "Azure Keyvault tenant ID")
)

func main() {
	flag.Parse()

	// keyvault name is required
	if len(*keyvaultName) == 0 {
		klog.Fatal("keyvault name not provided")
	}
	// keyvault secret name is required
	if len(*keyvaultSecretName) == 0 {
		klog.Fatal("keyvault secret name not provided")
	}
	// clientID and clientSecret are required to access secret from keyvault
	if len(*clientID) == 0 || len(*clientSecret) == 0 {
		klog.Fatal("clientID and client secret are required")
	}
	if len(*tenantID) == 0 {
		klog.Fatalf("tenantID is required")
	}

	// use azure public cloud environment for the purpose of demo
	env := azure.PublicCloud

	// create new keyvault client
	kvClient := kv.New()
	// get token to access keyvault
	tp, err := token(env, *clientID, *clientSecret, *tenantID)
	if err != nil {
		klog.Fatalf("failed to get token: %v", err)
	}
	kvClient.Authorizer = autorest.NewBearerAuthorizer(tp)

	vaultURL, err := getVaultURL(env, *keyvaultName)
	if err != nil {
		klog.Fatalf("failed to get vault url: %+v", err)
	}

	// get secret from keyvault
	secret, err := kvClient.GetSecret(context.Background(), *vaultURL, *keyvaultSecretName, *keyvaultSecretVersion)
	if err != nil {
		klog.Fatalf("failed to get secret from keyvault: %v", err)
	}
	// log the secret
	klog.Infof(*secret.Value)
}

func token(env azure.Environment, clientID, clientSecret, tenantID string) (*adal.ServicePrincipalToken, error) {
	kvEndPoint := env.KeyVaultEndpoint
	if '/' == kvEndPoint[len(kvEndPoint)-1] {
		kvEndPoint = kvEndPoint[:len(kvEndPoint)-1]
	}
	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return nil, err
	}
	return adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, kvEndPoint)
}

func getVaultURL(azureEnvironment azure.Environment, vaultName string) (vaultURL *string, err error) {
	// Key Vault name must be a 3-24 character string
	if len(vaultName) < 3 || len(vaultName) > 24 {
		return nil, fmt.Errorf("invalid vault name: %q, must be between 3 and 24 chars", vaultName)
	}
	// See docs for validation spec: https://docs.microsoft.com/en-us/azure/key-vault/about-keys-secrets-and-certificates#objects-identifiers-and-versioning
	isValid := regexp.MustCompile(`^[-A-Za-z0-9]+$`).MatchString
	if !isValid(vaultName) {
		return nil, fmt.Errorf("invalid vault name: %q, must match [-a-zA-Z0-9]{3,24}", vaultName)
	}

	vaultDNSSuffixValue := azureEnvironment.KeyVaultDNSSuffix
	vaultURI := "https://" + vaultName + "." + vaultDNSSuffixValue + "/"
	return &vaultURI, nil
}
