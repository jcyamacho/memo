package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jcyamacho/memo/cmd"
	"github.com/jcyamacho/memo/internal/fsstore"
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/jcyamacho/memo/internal/workspace"
)

const configDirEnv = "MEMO_CONFIG_DIR"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "memo: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dir, err := storeDir()
	if err != nil {
		return fmt.Errorf("resolve config directory: %w", err)
	}

	store, err := fsstore.New(dir)
	if err != nil {
		return fmt.Errorf("open memory store: %w", err)
	}

	service := memory.NewService(store, workspace.GitResolver{})
	return cmd.New(service, dir).ExecuteContext(ctx)
}

func storeDir() (string, error) {
	if dir := os.Getenv(configDirEnv); dir != "" {
		return filepath.Abs(dir)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "memo"), nil
}
