package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jcyamacho/memo/cmd"
	"github.com/jcyamacho/memo/internal/config"
	"github.com/jcyamacho/memo/internal/fsstore"
	"github.com/jcyamacho/memo/internal/memory"
	"github.com/jcyamacho/memo/internal/workspace"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "memo: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dir, err := config.Dir()
	if err != nil {
		return fmt.Errorf("resolve config directory: %w", err)
	}

	store, err := fsstore.New(dir)
	if err != nil {
		return fmt.Errorf("open memory store: %w", err)
	}

	service := memory.NewService(store, workspace.GitResolver{})
	return cmd.New(service).ExecuteContext(ctx)
}
