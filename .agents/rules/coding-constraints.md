# Coding Constraints

Applies to all Go code in this repository.

## Language and Dependencies

- Use Go 1.27.3 and the standalone modules
  `github.com/slaincsuite/modality-gopn-sdk-go` and
  `github.com/slaincsuite/modality-gopn-sdk-go/sample`.
- Runtime dependencies are limited to `hashicorp/go-plugin`,
  `hashicorp/go-hclog`, and `google.golang.org/grpc` unless the owner records a
  dated decision. Test-only dependencies must not become runtime imports.
- Never import applications, adapters, modules, domains, or packages from the
  modality monolith.

## Comments and Public API

- Match the Lark SDK layout style: add a short Chinese `//` doc line at the top
  of every hand-written code file. Do not add per-file license headers; the
  repository `LICENSE` governs the source.
- Every code file starts with one short Chinese comment. Additional comments
  explain non-obvious behavior only.
- Every exported identifier has a Chinese doc comment beginning with its name.
- Avoid package-level mutable state and side effects in `init()`.

## Errors

- Name known conditions `ErrXxx`; wrap with `%w`; make callers able to use
  `errors.Is`.
- Do not discard errors without a stated reason.

## Tests

- Prefer table-driven tests named `TestFunction_Scenario`.
- Use `testing`, `testify/require`, and explicit channel or context waits; never
  sleep-and-hope.
- Cover configuration validation and the plugin handshake before changing the
  public API.

## Verification

- Run `make verify` (root + sample) before claiming work is complete.
- Run tests with `-race` when changing lifecycle, concurrency, or process setup.
