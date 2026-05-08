# Changelog

## [Unreleased] - 2025-05-07

### Fixed
- Removed duplicate file `pkg/analyzer/kubernetes/kubernetes 2.go` that prevented compilation
- Fixed `notifier_test.go` field casing (`notifiers` → `Notifiers`)
- Fixed `history.go` non-unique ID generation using `crypto/rand`
- Fixed `pprof` global handler registration — pprof now uses a dedicated mux only when invoked
- Fixed online collector: added pagination (`Limit: 500`), removed deprecated `ComponentStatuses` API, replaced silent error swallowing with `klog` logging
- Fixed online collector thread safety with mutex-protected lazy client initialization
- Fixed online collector reliability with retry/backoff for K8s API calls
- Fixed `autofix` engine — `ConfirmationRequired` now blocks execution until confirmed
- Fixed `scanner` — exported `IsAvailable()` for external use
- Fixed AI provider — added timeout to `Analyze()` API call
- Fixed Operator logging — switched from development to production zap config
- Fixed Operator NodeDiagnostic Pod spec — added `HostNetwork`, `HostPID`, `Privileged`
- Fixed Operator ClusterDiagnostic — implemented Job result collection via ConfigMap
- Fixed Operator schedule controller — full cron expression support (`@monthly`, `@yearly`, `@every Nh`, 5-field cron)
- Fixed Dockerfile — non-root user, multi-stage build, proper LDFLAGS, safe COPY
- Fixed Makefile — added `linux-arm64`, `docker-build`, `docker-push` targets

### Added
- Wired 5 stub CLI commands: `fix`, `cost`, `scan`, `trace` (returns unimplemented error), `multicluster` (returns unimplemented error)
- Added `kudig ai` CLI command for AI-powered diagnosis
- Added eBPF probe user-facing warnings when BPF programs are unavailable
- Added tests for 6 previously untested packages: `servicemesh`, `autofix`, `cost`, `rca`, `scanner`, `tui` (28/28 packages now have tests)
- Added registry dispatch tests: `ExecuteAll`, `ExecuteByMode`, `ExecuteByCategory`, `ExecuteByNames`, dependency sorting, circular dependency detection, cancelled context handling
- Implemented TUI diagnosis pipeline — `startDiagnosis()` now calls collector → analyzer instead of returning empty results
- Converted Operator NodeDiagnostic from Job to DaemonSet for proper multi-node coverage
- DaemonSet uses node affinity for specific node targeting and tolerates all taints
- Per-node result collection via ConfigMaps (`kudig-result-{diagnostic}-{node}`)
- Created `docs/PRODUCTION_READINESS_AUDIT.md`
- Created `docs/FIX_RECORDS.md`

### Changed
- `kudig trace` and `kudig multicluster` now explicitly return errors instead of silently returning fake data
- Test coverage improved from ~30% to 52.5%

### Known Limitations
- eBPF analyzers (TCP/DNS/FileIO) are registered but return empty data — no BPF programs exist yet
- `kudig trace` and `kudig multicluster` are not implemented
- No E2E tests with real Kubernetes clusters
- CLI layer (`cmd/kudig/main.go`) has 0% test coverage
