package main

// This is the entrypoint: it wires config -> loadgen -> api together
// and starts listening. There's no test file for main() — the packages
// it calls are already tested individually; this file is verified by
// actually running the server (see docs/DOCKER.md and the quick checks
// in docs/ROADMAP.md's Phase 2 checklist, e.g. curl'ing /health).
//
// TODO, roughly in order:
//  1. cfg := config.Load()
//  2. mgr := loadgen.NewManager(cfg)
//  3. srv := api.NewServer(mgr)
//  4. log something like "listening on :<port>"
//  5. http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), srv.Routes())
//     and log.Fatal if it returns an error
//
// You'll need these imports once you fill this in:
//   "fmt", "log", "net/http",
//   "observability-demo/backend/internal/api",
//   "observability-demo/backend/internal/config",
//   "observability-demo/backend/internal/loadgen"

import (
	"fmt"
	"log"
	"net/http"

	"observability-demo/backend/internal/api"
	"observability-demo/backend/internal/config"
	"observability-demo/backend/internal/loadgen"
)

func main() {
	cfg := config.Load()
	mgr := loadgen.NewManager(cfg)
	srv := api.NewServer(mgr)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}
