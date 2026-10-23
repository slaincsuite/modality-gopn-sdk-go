# Host Access

- Every cross-module host read or controlled write is represented as a
  `host.Capability` with `read` or `controlled-write` mode.
- `host.Client` denies by default and cannot grant a capability absent from the
  host-provided declaration.
- `host.Authorization` is transport metadata, not proof. The host must call its
  authorization engine again on every request.
- Controlled writes are exceptional: they require an explicit host capability,
  per-call authorization, validation, and audit.
- The SDK never exposes host module, adapter, domain, or repository types.

## Extension Points

- `extensionpoint.Definition` belongs to the host and states singleton or
  multi-instance behavior.
- A plugin contributes only `extensionpoint.Extension` metadata and a service
  name already declared in `plugin.json`.
- Invocation order and enabling remain host decisions after authorization.
