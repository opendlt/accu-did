// Package main merges resolver and registrar OpenAPI specs into a unified SDK spec
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type OpenAPISpec struct {
	OpenAPI    string                 `yaml:"openapi"`
	Info       Info                   `yaml:"info"`
	Servers    []Server               `yaml:"servers,omitempty"`
	Paths      map[string]interface{} `yaml:"paths"`
	Components Components             `yaml:"components,omitempty"`
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version"`
}

type Server struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description,omitempty"`
}

type Components struct {
	Schemas map[string]interface{} `yaml:"schemas,omitempty"`
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: sdkmerge <resolver.yaml> <registrar.yaml> <output.yaml> [version]")
	}

	resolverPath := os.Args[1]
	registrarPath := os.Args[2]
	outputPath := os.Args[3]

	version := "1.0.0"
	if len(os.Args) > 4 {
		version = os.Args[4]
	}

	// Read resolver spec
	resolverData, err := os.ReadFile(resolverPath)
	if err != nil {
		log.Fatalf("Failed to read resolver spec: %v", err)
	}

	var resolverSpec OpenAPISpec
	if err := yaml.Unmarshal(resolverData, &resolverSpec); err != nil {
		log.Fatalf("Failed to parse resolver spec: %v", err)
	}

	// Read registrar spec
	registrarData, err := os.ReadFile(registrarPath)
	if err != nil {
		log.Fatalf("Failed to read registrar spec: %v", err)
	}

	var registrarSpec OpenAPISpec
	if err := yaml.Unmarshal(registrarData, &registrarSpec); err != nil {
		log.Fatalf("Failed to parse registrar spec: %v", err)
	}

	// Create merged spec
	merged := OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       "Accu-DID SDK API",
			Description: "Combined API specification for Accumulate DID resolver and registrar services",
			Version:     version,
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080",
				Description: "Resolver service",
			},
			{
				URL:         "http://localhost:8081",
				Description: "Registrar service",
			},
		},
		Paths:      make(map[string]interface{}),
		Components: Components{
			Schemas: make(map[string]interface{}),
		},
	}

	// Merge paths
	for path, pathSpec := range resolverSpec.Paths {
		merged.Paths[path] = pathSpec
	}
	for path, pathSpec := range registrarSpec.Paths {
		merged.Paths[path] = pathSpec
	}

	// Merge schemas with prefixes to avoid collisions
	if resolverSpec.Components.Schemas != nil {
		for name, schema := range resolverSpec.Components.Schemas {
			merged.Components.Schemas["Resolver_"+name] = schema
		}
	}
	if registrarSpec.Components.Schemas != nil {
		for name, schema := range registrarSpec.Components.Schemas {
			merged.Components.Schemas["Registrar_"+name] = schema
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write merged spec
	mergedData, err := yaml.Marshal(&merged)
	if err != nil {
		log.Fatalf("Failed to marshal merged spec: %v", err)
	}

	if err := os.WriteFile(outputPath, mergedData, 0644); err != nil {
		log.Fatalf("Failed to write merged spec: %v", err)
	}

	fmt.Printf("✅ Merged OpenAPI specs written to: %s\n", outputPath)
	fmt.Printf("   Resolver paths: %d\n", len(resolverSpec.Paths))
	fmt.Printf("   Registrar paths: %d\n", len(registrarSpec.Paths))
	fmt.Printf("   Total paths: %d\n", len(merged.Paths))
	fmt.Printf("   Total schemas: %d\n", len(merged.Components.Schemas))
}