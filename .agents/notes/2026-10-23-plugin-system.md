# 2026-10-23 - Plugin System Scope

## Decisions

1. The repository root is metadata-only. The former root package moved to
   `service`, matching the focused-package layout of the reference SDK.
2. `host` models explicit capabilities and default-deny authorization. A
   capability is read-only or controlled-write; the host still re-authorizes
   and audits every call.
3. `extensionpoint` models singleton and multi-instance extension points. A
   plugin declares an implementation and an already-declared gRPC service; it
   does not install host objects.
4. Manifest applications now include `grants` and `extensions`, allowing the
   host to pre-filter unauthorized or dangling declarations.
5. The SDK README is Chinese only and no changelog is maintained in this
   repository.
