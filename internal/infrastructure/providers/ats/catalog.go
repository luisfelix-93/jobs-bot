package ats

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CompanyCatalogEntry struct {
	Name            string   `yaml:"name"`
	Provider        string   `yaml:"-"` // Set dynamically based on YAML file name
	BoardToken      string   `yaml:"board_token"`
	Country         string   `yaml:"country"`
	RemoteFriendly  bool     `yaml:"remote_friendly"`
	Categories      []string `yaml:"categories"`
	CareerPageURL   string   `yaml:"career_page_url"`
	LastValidation  string   `yaml:"last_validation"`
}

type ProviderCatalog struct {
	Companies map[string]CompanyCatalogEntry `yaml:"companies"`
}

type CollectionsCatalog struct {
	Collections map[string][]string `yaml:"collections"`
}

type Catalog struct {
	Companies   map[string]CompanyCatalogEntry
	Collections map[string][]string
}

func LoadCatalog(catalogDir string) (*Catalog, error) {
	collectionsPath := filepath.Join(catalogDir, "collections.yaml")
	data, err := os.ReadFile(collectionsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collections: %w", err)
	}

	var colCat CollectionsCatalog
	if err := yaml.Unmarshal(data, &colCat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collections: %w", err)
	}

	entries, err := os.ReadDir(catalogDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read catalog directory: %w", err)
	}

	companies := make(map[string]CompanyCatalogEntry)
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "collections.yaml" || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		providerName := entry.Name()[:len(entry.Name())-len(filepath.Ext(entry.Name()))]
		filePath := filepath.Join(catalogDir, entry.Name())
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		var provCat ProviderCatalog
		if err := yaml.Unmarshal(fileData, &provCat); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %w", entry.Name(), err)
		}

		for id, company := range provCat.Companies {
			company.Provider = providerName
			companies[id] = company
		}
	}

	return &Catalog{
		Companies:   companies,
		Collections: colCat.Collections,
	}, nil
}

func (c *Catalog) ResolveCompanies(requestedCols []string, requestedComps []string) ([]CompanyCatalogEntry, error) {
	resolvedMap := make(map[string]CompanyCatalogEntry)

	// Resolve collections
	for _, colName := range requestedCols {
		comps, exists := c.Collections[colName]
		if !exists {
			return nil, fmt.Errorf("collection %q not found in catalog", colName)
		}
		for _, compID := range comps {
			comp, ok := c.Companies[compID]
			if !ok {
				return nil, fmt.Errorf("company %q in collection %q not found in catalog", compID, colName)
			}
			resolvedMap[compID] = comp
		}
	}

	// Resolve individual companies
	for _, compID := range requestedComps {
		comp, ok := c.Companies[compID]
		if !ok {
			return nil, fmt.Errorf("company %q not found in catalog", compID)
		}
		resolvedMap[compID] = comp
	}

	result := make([]CompanyCatalogEntry, 0, len(resolvedMap))
	for _, comp := range resolvedMap {
		result = append(result, comp)
	}
	return result, nil
}
