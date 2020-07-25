package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

// Context is struct for parsing context section of kubeconfig file
type Context struct {
	Context struct {
		Namespace string `yaml:"namespace"`
	} `yaml:"context"`
	Name string `yaml:"name"`
}

// Kubeconfig is struct for parsing kubeconfig file
type Kubeconfig struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

func getCurrentContext(configs []Kubeconfig) string {
	for _, config := range configs {
		if config.CurrentContext != "" {
			return config.CurrentContext
		}
	}
	return ""
}

func getContextNamespace(context string, configs []Kubeconfig) string {
	for _, config := range configs {
		for _, ctx := range config.Contexts {
			if ctx.Name == context {
				return ctx.Context.Namespace
			}
		}
	}
	return "default"
}

var allowedOutput = []string{"json", "slug", "context", "namespace"}

func validateOutputFlag(output string) string {
	for _, variant := range allowedOutput {
		if output == variant {
			return ""
		}
	}

	return fmt.Sprintf(`Unexpected output value "%s". Allowed values: %s`, output, strings.Join(allowedOutput, ", "))
}

func main() {
	kubeconfigRawEnv := os.Getenv("KUBECONFIG")
	if kubeconfigRawEnv == "" {
		fmt.Fprintln(os.Stderr, "Empty KUBECONFIG env variable")
		os.Exit(1)
	}

	outputFlag := flag.String("o", "slug", `Output values. Either "json", "context", "namespace" or "slug" (context/namespace)`)
	separatorFlag := flag.String("s", "/", `Separator for slug. Create stylish prompt with -s='⎈ '`)
	flag.Parse()

	validationError := validateOutputFlag(*outputFlag)
	if validationError != "" {
		fmt.Println(validationError)
		os.Exit(1)
	}

	output := *outputFlag
	separator := *separatorFlag

	kubeconfigFilePaths := strings.Split(kubeconfigRawEnv, ":")
	for i, path := range kubeconfigFilePaths {
		kubeconfigFilePaths[i] = strings.TrimSpace(path)
	}

	parsedKubeconfigs := make([]Kubeconfig, len(kubeconfigFilePaths))

	var wg sync.WaitGroup
	wg.Add(len(kubeconfigFilePaths))

	for i, file := range kubeconfigFilePaths {
		// Parse all files, skip files with syntax errors
		go func(fileName string, i int) {
			defer wg.Done()

			var config Kubeconfig
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file %s %v", fileName, err)
				return
			}

			err = yaml.NewDecoder(file).Decode(&config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %s %v", fileName, err)
				return
			}

			file.Close()
			parsedKubeconfigs[i] = config
		}(file, i)
	}

	wg.Wait()

	currentContext := getCurrentContext(parsedKubeconfigs)

	if output == "context" {
		fmt.Print(currentContext)
	}

	currentNamespace := getContextNamespace(currentContext, parsedKubeconfigs)

	if output == "namespace" {
		fmt.Print(currentNamespace)
	}

	if output == "slug" {
		fmt.Printf("%s%s%s", currentContext, separator, currentNamespace)
	}

	if output == "json" {
		jsonOutput := struct {
			Context   string `json:"context"`
			Namespace string `json:"namespace"`
		}{
			Context:   currentContext,
			Namespace: currentNamespace,
		}

		bytes, err := json.Marshal(&jsonOutput)
		if err != nil {
			fmt.Printf("Error while formatting json\n%s", err)
			os.Exit(1)
		}

		fmt.Print(string(bytes))
	}
}
