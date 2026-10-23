# Naming Conventions

Applies to directories, files, APIs, and services.

## Files and Packages

- Package names are lowercase, singular where possible, and describe what they
  provide; this repository exposes `modality`, `core`, and `manifest`.
- Go file names are lowercase and short; only `_test.go` may contain an
  underscore for its standard suffix.
- README is `README.md`; one per public package.

## Identifiers

- Export names describe plugin concepts: `Plugin`, `Service`, `RegisterFunc`.
  Do not use generic names such as `util`, `helper`, `base`, or `manager`.
- Business service names follow the proto package, for example
  `greeter.v1.Greeter`, and must match `plugin.json`.
- Keep service names stable after release; use a new name for an incompatible
  RPC contract.
