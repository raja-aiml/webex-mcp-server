package app

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/raja-aiml/webex-mcp-server-go/internal/testutil"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Name:     "test-app",
		Version:  "1.0.0",
		HTTPAddr: ":8080",
	}

	app := New(cfg)

	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.config.Name != cfg.Name {
		t.Errorf("app.config.Name = %v, want %v", app.config.Name, cfg.Name)
	}

	if app.config.Version != cfg.Version {
		t.Errorf("app.config.Version = %v, want %v", app.config.Version, cfg.Version)
	}

	if app.config.HTTPAddr != cfg.HTTPAddr {
		t.Errorf("app.config.HTTPAddr = %v, want %v", app.config.HTTPAddr, cfg.HTTPAddr)
	}

	if app.ctx == nil {
		t.Error("app.ctx is nil")
	}

	if app.cancel == nil {
		t.Error("app.cancel is nil")
	}
}

func TestApp_Shutdown(t *testing.T) {
	app := New(Config{
		Name:     "test-app",
		Version:  "1.0.0",
		HTTPAddr: "",
	})

	// Create a channel to verify context cancellation
	ctxDone := make(chan struct{})
	go func() {
		<-app.ctx.Done()
		close(ctxDone)
	}()

	// Call shutdown
	err := app.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Verify context was cancelled
	select {
	case <-ctxDone:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("Context was not cancelled within timeout")
	}
}

func TestApp_Run_MissingAPIKey(t *testing.T) {
	// Ensure API key is not set
	os.Unsetenv("WEBEX_PUBLIC_WORKSPACE_API_KEY")

	app := New(Config{
		Name:     "test-app",
		Version:  "1.0.0",
		HTTPAddr: "",
	})

	err := app.Run()
	if err == nil {
		t.Error("Expected error for missing API key, got nil")
	}
}

func TestApp_Run_SignalHandling(t *testing.T) {
	// Skip in short mode as this test involves signals
	if testing.Short() {
		t.Skip("Skipping signal handling test in short mode")
	}

	// Set required API key
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer cleanup()

	app := New(Config{
		Name:     "test-app",
		Version:  "1.0.0",
		HTTPAddr: "",
	})

	// Run app in goroutine
	runErr := make(chan error, 1)
	go func() {
		runErr <- app.Run()
	}()

	// Give app time to start
	time.Sleep(100 * time.Millisecond)

	// Send interrupt signal to self
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	
	if err := proc.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}

	// Wait for app to exit
	select {
	case err := <-runErr:
		if err != nil {
			t.Errorf("Run() returned error on signal: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("App did not exit within timeout after signal")
	}
}

func TestApp_ContextCancellation(t *testing.T) {
	app := New(Config{
		Name:     "test-app",
		Version:  "1.0.0",
		HTTPAddr: "",
	})

	// Verify initial context is not done
	select {
	case <-app.ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Success
	}

	// Cancel context
	app.cancel()

	// Verify context is done
	select {
	case <-app.ctx.Done():
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Context should be done after cancel")
	}
}