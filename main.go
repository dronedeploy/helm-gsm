package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	flag "github.com/spf13/pflag"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"gopkg.in/yaml.v3"
)

const (
	gsmPrefix = "gsm:"
)

var (
	b64Encode = false
)

func parseNode(ctx context.Context, sc *secretmanager.Client, node interface{}) interface{} {
	switch node.(type) {
	case map[string]interface{}:
		for k, v := range node.(map[string]interface{}) {
			node.(map[string]interface{})[k] = parseNode(ctx, sc, v)
		}
		return node
	case []interface{}:
		for i := range node.([]interface{}) {
			node.([]interface{})[i] = parseNode(ctx, sc, node.([]interface{})[i])
		}
		return node
	case string:
		if strings.HasPrefix(node.(string), gsmPrefix) {
			checkValidSecret(node.(string))
			secretPath := strings.TrimPrefix(node.(string), gsmPrefix)
			if b64Encode == true {
				node = base64.StdEncoding.EncodeToString(getSecret(ctx, sc, transformStringToCanonicalName(secretPath)))
				return node
			}
			node = string(getSecret(ctx, sc, transformStringToCanonicalName(secretPath)))
		}
		return node
	default:
		return node
	}
}

func transformStringToCanonicalName(raw string) string {
	s := strings.Split(raw, "/")
	return fmt.Sprintf("projects/%s/secrets/%s/versions/%s", s[0], s[1], s[2])
}

func getSecret(ctx context.Context, sc *secretmanager.Client, s string) []byte {
	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: s,
	}

	// Call the API.
	result, err := sc.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	return result.Payload.Data
}

func checkValidSecret(s string) bool {
	regex_string := `^gsm:[a-z][a-z0-9-]{4,28}[a-z0-9]\/[a-zA-Z0-9_-]+\/(latest|[1-9][0-9]*)$`
	matched, err := regexp.Match(regex_string, []byte(s))
	if err != nil {
		log.Fatalf("error matching regex")
	}
	if !matched {
		log.Fatal("Secret did not match required format.\n" +
			"\n" +
			"Must be in the form 'gsm:project_id/secret_name/version'.\n" +
			"project_id: 'The unique, user-assigned ID of the Project. It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited.'\n" +
			"secret_name: 'Secret names can only contain English letters (A-Z), numbers (0-9), dashes (-), and underscores (_)'\n" +
			"version: 'Versions are a monotonically increasing integer starting at 1.'\n" +
			"\n" +
			"example: 'gsm:project-id/secret_name/1'\n" +
			"regex: '" + regex_string + "'\n" +
			"\n")
	}
	return matched
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	ctx := context.Background()

	// Create the secretmanager client once and reuse it throughout
	sc, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to setup SecretManager client: %v", err)
	}

	// Grab the input filename, abort if it doesn't exist or is a directory.
	var inputFilename string
	flag.StringVarP(&inputFilename, "secretsFile", "f", "secrets.yaml", "Filepath to a yaml file with secrets")
	flag.BoolVarP(&b64Encode, "b64Encode", "b", false, "base64 encode values in the output .dec file")
	flag.Parse()

	fmt.Println(inputFilename)
	if !fileExists(inputFilename) {
		log.Fatalf("input file %v does not exist or is a directory\n", inputFilename)
	}

	// Unmarshal the input file
	inputFile, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		log.Printf("failed to read input file %v ", err)
	}
	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(inputFile, &data)
	if err != nil {
		log.Fatalf("failed unmarshalling input data: %v", err)
	}

	// Loop over the nodes of the YAML parsing them recursively
	for k, v := range data {
		data[k] = parseNode(ctx, sc, v)
	}

	// Marshal the parsed data back to YAML and write to the output file
	outputYAML, err := yaml.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	outputFilename := inputFilename + ".dec"
	log.Printf("Writing plaintext secrets to %s.", outputFilename)
	err = ioutil.WriteFile(outputFilename, outputYAML, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
