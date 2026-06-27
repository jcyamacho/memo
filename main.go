package main

import (
	"context"
	"log"
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
	dir, err := config.Dir()
	if err != nil {
		log.Fatalf("resolve config directory: %v", err)
	}

	store, err := fsstore.New(dir)
	if err != nil {
		log.Fatalf("open memory store: %v", err)
	}
	cmd.SetService(memory.NewService(store, workspace.GitResolver{}))

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd.Execute(ctx)
}
