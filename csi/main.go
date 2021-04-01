package main

import (
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
	"k8s.io/klog/v2"
)

var (
	secretsPath = flag.String("secrets-path", "/mnt/secrets-store", "fs path that contains the secrets")
	secretName  = flag.String("secret-name", "app-secret", "secret name")

	watcher *fsnotify.Watcher
	secretValue string
)

func main() {
	flag.Parse()

	if len(*secretName) == 0 {
		klog.Fatal("secret name not provided")
	}

	secretFile := filepath.Join(*secretsPath, *secretName)
	secret, err := getSecret(secretFile)
	if err != nil {
		klog.Fatalf("failed to get secret file: %v", err)
	}
	secretValue = secret
	// log the secret
	klog.Infof(secret)

	// To receive automatic updates and refresh when the secret in file system is rotated
	// we will setup a file watcher using fsnotify
	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// add the secrets dir to the watcher
	if err = watcher.Add(*secretsPath); err != nil {
		klog.Fatalf("failed to add secret path to file watcher: %v", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case <-watcher.Events:
				secret, err := getSecret(secretFile)
				if err != nil {
					klog.Fatalf("failed to get secret file: %v", err)
				}
				if secretValue != secret {
					secretValue = secret
					// log the secret
					klog.Infof(secret)
				}

				// watch for errors
			case err := <-watcher.Errors:
				klog.Error("failed to watch: %v", err)
			}
		}
	}()

	<-done
}

// getSecret reads the secret from the pod file system
func getSecret(secretFile string) (string, error) {
	// read the secret from pod file system
	if _, err := os.Stat(secretFile); err != nil {
		return "", err
	}
	bytes, err := os.ReadFile(secretFile)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
