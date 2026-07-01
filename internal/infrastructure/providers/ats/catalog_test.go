package ats

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCatalogAndResolve(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "catalog-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	collectionsData := `
collections:
  tech:
    - vercel
    - stripe
  ai:
    - openai
`
	greenhouseData := `
companies:
  openai:
    name: "OpenAI"
    board_token: "openai"
    country: "US"
    remote_friendly: true
  stripe:
    name: "Stripe"
    board_token: "stripe"
    country: "US"
    remote_friendly: true
`
	leverData := `
companies:
  vercel:
    name: "Vercel"
    board_token: "vercel"
    country: "US"
    remote_friendly: true
`

	if err := os.WriteFile(filepath.Join(tempDir, "collections.yaml"), []byte(collectionsData), 0644); err != nil {
		t.Fatalf("failed to write collections.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "greenhouse.yaml"), []byte(greenhouseData), 0644); err != nil {
		t.Fatalf("failed to write greenhouse.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "lever.yaml"), []byte(leverData), 0644); err != nil {
		t.Fatalf("failed to write lever.yaml: %v", err)
	}

	cat, err := LoadCatalog(tempDir)
	if err != nil {
		t.Fatalf("failed to load catalog: %v", err)
	}

	if len(cat.Collections) != 2 {
		t.Errorf("expected 2 collections, got %d", len(cat.Collections))
	}
	if len(cat.Companies) != 3 {
		t.Errorf("expected 3 companies, got %d", len(cat.Companies))
	}

	// Verify provider was set from filename
	openai, ok := cat.Companies["openai"]
	if !ok || openai.Provider != "greenhouse" {
		t.Errorf("expected openai provider to be 'greenhouse', got %q", openai.Provider)
	}

	vercel, ok := cat.Companies["vercel"]
	if !ok || vercel.Provider != "lever" {
		t.Errorf("expected vercel provider to be 'lever', got %q", vercel.Provider)
	}

	// Test ResolveCompanies
	resolved, err := cat.ResolveCompanies([]string{"ai"}, []string{"stripe"})
	if err != nil {
		t.Fatalf("failed to resolve companies: %v", err)
	}

	if len(resolved) != 2 {
		t.Errorf("expected 2 resolved companies, got %d", len(resolved))
	}

	// Make sure resolved entries contain openai and stripe
	foundOpenAI := false
	foundStripe := false
	for _, comp := range resolved {
		if comp.BoardToken == "openai" {
			foundOpenAI = true
		}
		if comp.BoardToken == "stripe" {
			foundStripe = true
		}
	}

	if !foundOpenAI {
		t.Error("resolved list missing OpenAI")
	}
	if !foundStripe {
		t.Error("resolved list missing Stripe")
	}

	// Test error cases
	_, err = cat.ResolveCompanies([]string{"invalid-col"}, nil)
	if err == nil {
		t.Error("expected error for invalid collection name, got nil")
	}

	_, err = cat.ResolveCompanies(nil, []string{"invalid-comp"})
	if err == nil {
		t.Error("expected error for invalid company name, got nil")
	}
}
