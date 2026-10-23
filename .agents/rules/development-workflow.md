# Development Workflow

Follow this order for every non-trivial change.

1. Inventory the root package, `core/`, `manifest/`, `sample/`, this index, and
   the notes directory before creating a file;
   extend an existing artifact when its responsibility already exists.
2. Keep the public API narrow. Add options only when a caller needs them.
3. Preserve the host contract described in
   [release-compatibility](release-compatibility.md).
4. Add or update focused tests for validation and lifecycle changes.
5. Run `make verify`; use `make test-race` for process or concurrency changes.
6. Record durable decisions in a dated `.agents/notes` file and update the
   notes index in `AGENTS.md`.
7. Never commit unless the owner explicitly requests it.
