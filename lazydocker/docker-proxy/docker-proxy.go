// #ddev-generated
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

const (
	dockerSocket = "/var/run/docker-host.sock"
	proxySocket  = "/var/run/docker.sock"
)

func getMatchString() string {
	project := os.Getenv("DDEV_PROJECT")
	if project == "" {
		log.Fatal("❌ DDEV_PROJECT environment variable is not set")
	}
	return "ddev-" + project
}

func main() {
	matchString := getMatchString()

	if _, err := os.Stat(proxySocket); err == nil {
		_ = os.Remove(proxySocket)
	}

	listener, err := net.Listen("unix", proxySocket)
	if err != nil {
		log.Fatalf("❌ Failed to listen on socket: %v", err)
	}
	defer listener.Close()

	log.Printf("🚀 Docker proxy started on %s (filter: %s)", proxySocket, matchString)

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "docker"
		},
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", dockerSocket)
			},
		},
		ModifyResponse: func(resp *http.Response) error {
			path := resp.Request.URL.Path
			log.Printf("➡️  Proxy intercepted: %s", path)

			cleanPath := stripVersionPrefix(path)

			switch {
			case cleanPath == "/containers/json" && resp.StatusCode == 200:
				return filterContainersList(resp, matchString)

			case cleanPath == "/images/json" && resp.StatusCode == 200:
				return filterImages(resp, matchString)

			case cleanPath == "/volumes":
				return filterVolumes(resp, matchString)

			case cleanPath == "/networks":
				return filterNetworks(resp, matchString)
			}

			return nil
		},
	}

	server := &http.Server{Handler: proxy}
	log.Fatal(server.Serve(listener))
}

func stripVersionPrefix(path string) string {
	if strings.HasPrefix(path, "/v") {
		parts := strings.SplitN(path, "/", 3)
		if len(parts) >= 3 {
			return "/" + parts[2]
		}
	}
	return path
}

func filterContainersList(resp *http.Response, match string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var containers []map[string]interface{}
	if err := json.Unmarshal(body, &containers); err != nil {
		return err
	}

	filtered := []map[string]interface{}{}
	for _, c := range containers {
		if names, ok := c["Names"].([]interface{}); ok {
			for _, n := range names {
				if strings.Contains(n.(string), match) {
					filtered = append(filtered, c)
					break
				}
			}
		}
	}

	log.Printf("🔍 Filtered containers: kept %d of %d", len(filtered), len(containers))

	newBody, err := json.Marshal(filtered)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(strings.NewReader(string(newBody)))
	resp.ContentLength = int64(len(newBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
	resp.Header.Set("Content-Type", "application/json")
	return nil
}

func filterContainerInspectByName(resp *http.Response, match string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var container map[string]interface{}
	if err := json.Unmarshal(body, &container); err != nil {
		return err
	}

	name, ok := container["Name"].(string)
	if !ok || !strings.Contains(name, match) {
		log.Printf("⛔ Blocking container inspect: %s", name)
		blockedBody := `{"message":"container not found"}`
		resp.Body = io.NopCloser(strings.NewReader(blockedBody))
		resp.ContentLength = int64(len(blockedBody))
		resp.Header.Set("Content-Length", strconv.Itoa(len(blockedBody)))
		resp.Header.Set("Content-Type", "application/json")
		resp.StatusCode = 404
		resp.Status = "404 Not Found"
	}
	return nil
}

func filterImages(resp *http.Response, match string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var images []map[string]interface{}
	if err := json.Unmarshal(body, &images); err != nil {
		return err
	}

	filtered := []map[string]interface{}{}
	for _, img := range images {
		if tags, ok := img["RepoTags"].([]interface{}); ok {
			for _, tag := range tags {
				if tagStr, ok := tag.(string); ok && strings.Contains(tagStr, match) {
					filtered = append(filtered, img)
					break
				}
			}
		}
	}

	log.Printf("🧊 Filtered images: kept %d of %d", len(filtered), len(images))

	newBody, err := json.Marshal(filtered)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(strings.NewReader(string(newBody)))
	resp.ContentLength = int64(len(newBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
	resp.Header.Set("Content-Type", "application/json")
	return nil
}

func filterVolumes(resp *http.Response, match string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	volumes, ok := data["Volumes"].([]interface{})
	if !ok {
		return nil // let it through
	}

	filtered := []interface{}{}
	for _, vol := range volumes {
		volMap, ok := vol.(map[string]interface{})
		if ok && strings.Contains(volMap["Name"].(string), match) {
			filtered = append(filtered, volMap)
		}
	}
	data["Volumes"] = filtered

	newBody, _ := json.Marshal(data)
	resp.Body = io.NopCloser(strings.NewReader(string(newBody)))
	resp.ContentLength = int64(len(newBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
	resp.Header.Set("Content-Type", "application/json")
	return nil
}

func filterNetworks(resp *http.Response, match string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var networks []map[string]interface{}
	if err := json.Unmarshal(body, &networks); err != nil {
		return err
	}

	filtered := []map[string]interface{}{}
	for _, net := range networks {
		if name, ok := net["Name"].(string); ok && strings.Contains(name, match) {
			filtered = append(filtered, net)
		}
	}

	newBody, _ := json.Marshal(filtered)
	resp.Body = io.NopCloser(strings.NewReader(string(newBody)))
	resp.ContentLength = int64(len(newBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
	resp.Header.Set("Content-Type", "application/json")
	return nil
}

