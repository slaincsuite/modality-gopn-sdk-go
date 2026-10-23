# 2026-10-23 - Plugin SDK Bootstrap

Context: the modality monolith already contains `library/gopn`, its host-side
runner. The owner asked for a standalone SDK repository for plugin developers,
including `AGENTS.md` and `.agents`. An initial thin `sdk/` implementation was
rejected because it invented an import path and required a non-existent
`RegisterGRPC` method. This note records the rebuilt contract.

## Decisions

1. **Standalone module.** The module path is
   `github.com/slaincsuite/modality-gopn-sdk-go`; it must not import modality
   application packages. This keeps plugin builds independent of host internals.
2. **Contract parity.** Handshake values are copied from modality
   `library/gopn/runner`: protocol `1`, cookie key `GOPN_PLUGIN`, and cookie
   value `gopn-plugin-handshake`. The SDK never renames the host contract.
3. **Lark SDK layout.** The repository uses a root public package, focused
   sub-packages, an independent `sample/` module, bilingual docs, changelog,
   Makefile, and CI, mirroring `larksuite/oapi-sdk-go` rather than a single
   `sdk/` directory.
4. **Generated-code bridge.** `Service.Register` is a
   `func(grpc.ServiceRegistrar)`. Plugin authors call the generated
   `RegisterXxxServer` inside the closure, so there is no invented interface.
5. **Single plugin key.** The go-plugin `Plugins` map contains one SDK-owned
   entry, `modality`; it registers every declared gRPC service. This avoids
   duplicate registration and lets the host use its own service factories.
6. **Package validation.** `manifest` mirrors the host rules for manifest,
   binary, UI, routes, discovery, and host version constraints so packaging
   failures happen before launch.
7. **Optional broker hook.** `WithBroker` forwards go-plugin's
   `GRPCBroker`; business callback protos remain a host/plugin contract.

## Current State

- Root package exposes `Plugin`, `Service`, `NewService`, `Serve`,
  `MustServe`, and functional options; duplicate or invalid services are
  rejected before serving.
- `core` owns handshake constants, config validation, and leveled logging.
- `manifest` loads and validates single packages and directories, checks
  SemVer ranges, and aggregates field errors.
- `sample` contains a real proto, generated Go code, a runnable plugin, and a
  subprocess E2E test that performs `Greeter/Hello` over gRPC.
- `make verify` is the single verification entry point.

## Open Items

- Add a broker E2E test once the host defines a callback proto and broker ID
  allocation strategy.
- Consider a generated helper for hosts to map `plugin.json services` to typed
  client factories.
- Publish release assets for the sample only after a versioning policy exists.
