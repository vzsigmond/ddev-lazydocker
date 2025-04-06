//#ddev-generated
package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	project := os.Getenv("DDEV_PROJECT")
	if project == "" {
		log.Fatal("❌ DDEV_PROJECT not set")
	}

	log.Printf("📦 Starting scoped LazyDocker for DDEV project: %s", project)

	// Start the docker proxy in the background
	go func() {
		cmd := exec.Command("/bin/docker-proxy")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Start(); err != nil {
			log.Fatalf("❌ Failed to start proxy: %v", err)
		}
	}()

	// Wait a second for proxy to initialize
	time.Sleep(1 * time.Second)

	// Start lazydocker with proxy
	cmd := exec.Command("lazydocker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Failed to start lazydocker: %v", err)
	}
}
