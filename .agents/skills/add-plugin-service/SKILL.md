---
name: add-plugin-service
description: Add or update a gRPC service exposed by a modality plugin using this SDK.
---

1. Confirm the proto contract exists and the generated code exports
   `RegisterXxxServer(grpc.ServiceRegistrar, XxxServer)`.
2. Inspect root `Plugin` / `Service` and `sample/cmd/greeter-plugin/main.go`
   before adding API; extend existing types instead of creating parallel APIs.
3. Choose a stable service name matching the proto full name and `plugin.json`.
4. Register it with `modality.NewService`, calling the generated register
   function inside the closure.
5. Add a subprocess test in `sample/` when registration or serving changes.
6. Run `make verify`; use `make test-race` for lifecycle changes.
7. Update `README.md`, `README.zh.md`, and `changelog.md` for public changes.
