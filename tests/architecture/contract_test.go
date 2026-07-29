package architecture_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRequiredArchitectureContract(t *testing.T) {
	t.Parallel()

	repoRoot := filepath.Join("..", "..")

	appcheckBytes, err := os.ReadFile(filepath.Join(repoRoot, "appcheck.toml"))
	if err != nil {
		t.Fatalf("read appcheck.toml: %v", err)
	}

	appcheck := string(appcheckBytes)
	for _, required := range []string{
		`required_checks = ["architecture"]`,
		"[architecture]",
		"enabled = true",
		"required = true",
		`baseline_file = ".appcheck/architecture-baseline.json"`,
	} {
		if !strings.Contains(appcheck, required) {
			t.Errorf("appcheck.toml missing required architecture declaration %q", required)
		}
	}

	if _, err := os.Stat(filepath.Join(repoRoot, "docs", "architecture", "ARCHITECTURE.md")); err != nil {
		t.Errorf("architecture document missing: %v", err)
	}

	baselineBytes, err := os.ReadFile(filepath.Join(repoRoot, ".appcheck", "architecture-baseline.json"))
	if err != nil {
		t.Fatalf("read architecture baseline: %v", err)
	}

	var baseline struct {
		Violations map[string]int `json:"violations"`
	}
	if err := json.Unmarshal(baselineBytes, &baseline); err != nil {
		t.Fatalf("parse architecture baseline: %v", err)
	}

	makefileBytes, err := os.ReadFile(filepath.Join(repoRoot, "Makefile"))
	if err != nil {
		t.Fatalf("read Makefile: %v", err)
	}

	makefile := string(makefileBytes)
	for _, required := range []string{"architecture-check:", "test-v1:", "test-v2:"} {
		if !strings.Contains(makefile, required) {
			t.Errorf("Makefile missing required gate %q", required)
		}
	}
}
