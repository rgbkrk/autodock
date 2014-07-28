package main

import (
	"encoding/json"

	"log"
	"os"
	"os/exec"
	"strings"

	"net/http"
)

var repoCommands = make(map[string][]string)

// DockerHubPayload is (partially) the webhook produced when a build completes on hub.docker.com
type DockerHubPayload struct {
	CallbackURL string              `json:"callback_url,omitempty"`
	Repository  DockerHubRepository `json:"repository,omitempty"`
	// There are a lot more fields, only putting in the ones I needed right now
	// Submit a PR if you'd like to see more
}

// DockerHubRepository Some of the pertinent bits that we'll use
type DockerHubRepository struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Owner     string `json:"owner,omitempty"`
	RepoName  string `json:"repo_name,omitempty"`
}

// DockerHubHandler accepts the webhook payload as produced when a build completes on hub.docker.com
func DockerHubHandler(rw http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var dockerHubPayload DockerHubPayload
	err := decoder.Decode(&dockerHubPayload)

	if err != nil {
		log.Println("Unable to read payload as JSON")
		log.Println(err)
		return
	}

	repo := dockerHubPayload.Repository.RepoName
	command, ok := repoCommands[repo]

	if !ok { // || len(command) == 0 {
		log.Printf("Repository \"%s\" not enabled\n", repo)
		return
	}

	log.Println("Processing", repo)
	log.Println("Running", command)

	repoCmd := exec.Command(command[0], command[1:]...)

	output, err := repoCmd.Output()
	if err != nil {
		panic(err)
	}

	log.Println(string(output))

}

func main() {

	// Accept all environment variables starting with AUTODOCK_ as configuration
	// for autodock. We're doing this instead of using flags because of how the
	// google/golang-runtime image works (no command line arguments).
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		key := pair[0]
		if strings.Contains(key, "AUTODOCK_") {
			autodockit := os.Getenv(key)
			autodocksplit := strings.Split(autodockit, ":")

			repo := autodocksplit[0]
			commands := strings.Join(autodocksplit[1:], ":")

			// Lazy mode, accepting spaces as delimiters between command arguments
			repoCommands[repo] = strings.Split(commands, " ")

		}
	}

	if len(repoCommands) == 0 {
		panic("No repositories configured")
	}

	log.Println("Docker repository actions:")
	for repo, commands := range repoCommands {
		log.Printf("\t%s: %s\n", repo, commands)
	}

	// TODO: Print out the full URL for the webhook on registry.hub.docker.com
	// TODO: Determine IP of the public host?
	log.Println("Point your Hook config at: http://{IP+Port}/autodock/v1/")

	http.HandleFunc("/autodock/v1/", DockerHubHandler)
	http.ListenAndServe(":8080", nil)
}
