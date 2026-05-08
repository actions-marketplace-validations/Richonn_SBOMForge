package sign

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/Richonn/sbomforge/internal/config"
)

func fakeExecCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	for _, a := range args {
		if strings.HasPrefix(a, "--bundle=") {
			bundlePath := strings.TrimPrefix(a, "--bundle=")
			_ = os.WriteFile(bundlePath, []byte(`{}`), 0644)
		}
	}
	os.Exit(0)
}

func TestSign_Disabled(t *testing.T) {
	cfg := &config.Config{Sign: false}

	bundlePath, err := Sign(context.Background(), cfg, "/tmp/sbom.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bundlePath != "" {
		t.Errorf("expected empty bundlePath when Sign=false, got %q", bundlePath)
	}
}

func TestSign_Success(t *testing.T) {
	cfg := &config.Config{Sign: true}

	origExecCommand := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = origExecCommand }()

	sbomFile, err := os.CreateTemp(t.TempDir(), "sbom-*.json")
	if err != nil {
		t.Fatal(err)
	}
	sbomFile.Close()

	bundlePath, err := Sign(context.Background(), cfg, sbomFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bundlePath != sbomFile.Name()+".bundle" {
		t.Errorf("unexpected bundlePath: %s", bundlePath)
	}
}

func TestSign_CosignFailure(t *testing.T) {
	cfg := &config.Config{Sign: true}

	execCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "false")
	}
	defer func() {
		execCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
			return exec.CommandContext(ctx, name, args...)
		}
	}()

	_, err := Sign(context.Background(), cfg, "/tmp/sbom.json")
	if err == nil {
		t.Error("expected error when cosign fails, got nil")
	}
}
