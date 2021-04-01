package main

import (
	"context"
	"flag"
	"fmt"
	"regexp"
	"time"

	kv "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest/azure"
	"k8s.io/klog/v2"
)

var (
	keyvaultName          = flag.String("keyvault-name", "kubecon-eu-2021", "Azure Keyvault name")
	keyvaultSecretName    = flag.String("keyvault-secret-name", "app-secret", "Azure Keyvault secret name")
	keyvaultSecretVersion = flag.String("keyvault-secret-version", "", "Azure Keyvault secret version")
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

	// use azure public cloud environment for the purpose of demo
	env := azure.PublicCloud

	// create new keyvault client
	kvClient := kv.New()
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		klog.Fatalf("failed to create authorizer: %v",err)
	}
	kvClient.Authorizer = authorizer

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
	// sleep forever
	time.Sleep(3600 * time.Second)
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
