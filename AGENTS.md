# AGENTS.md

Entry point for AI agents working in `modality-gopn-sdk-go`. Read this file
first, then load only the rules and skills relevant to the current task.

| Section | Path |
| --- | --- |
| Rules (must follow) | [.agents/rules/](.agents/rules/) |
| Skills (recurring procedures) | [.agents/skills/](.agents/skills/) |
| Notes (session records) | [.agents/notes/](.agents/notes/) |

## Project Overview

`modality-gopn-sdk-go` is the public Go SDK for third-party modality plugin
authors. The root directory stays metadata-only; all implementation lives in
focused packages following the SDK-side package-per-capability organization
expected by modality plugin authors.

| Aspect | Value |
| --- | --- |
| Module | `github.com/slaincsuite/modality-gopn-sdk-go` |
| Go | 1.27.3 |
| Transport | `hashicorp/go-plugin`, gRPC only |
| Public layout | `service`, `core`, `manifest`, `extensionpoint`, `host` |
| Host counterpart | modality `library/gopn` |

## Technology Stack

| Concern | Choice |
| --- | --- |
| Plugin transport | `hashicorp/go-plugin` v1.8.0 |
| RPC | `google.golang.org/grpc` v1.82.0 |
| License | GPL-3.0-only |
| Logging bridge | `hashicorp/go-hclog` |
| Manifest versions | `Masterminds/semver/v3` |
| Testing | stdlib `testing`, testify in examples |

## Directory Map

```text
modality-gopn-sdk-go/
├── AGENTS.md / README.md / Makefile / go.mod
├── service/          # Plugin, Service, Serve, gRPC bridge, logger adapter
├── core/             # handshake, config, logger, SDK version
├── manifest/         # plugin.json / action.json / routes.json models and validation
├── extensionpoint/   # plugin-side extension declarations
├── host/             # capability authorization and authorized host calls
└── .agents/          # repository-local agent knowledge
```

There are no Go source files in the repository root.

## Dependency Direction

```text
service ──▶ core
        ──▶ extensionpoint / host (public type re-exports only)
manifest ──▶ core-compatible data models only
```

- No SDK package imports a modality host application, adapter, module, domain,
  or internal package.
- `core` is the lowest-level package and must not import sibling packages.
- `manifest`, `extensionpoint`, and `host` must not import `service`.

## Rules Index

| File | Scope |
| --- | --- |
| [coding-constraints.md](.agents/rules/coding-constraints.md) | Go API and test style |
| [naming-conventions.md](.agents/rules/naming-conventions.md) | package, service, file naming |
| [development-workflow.md](.agents/rules/development-workflow.md) | inventory, verify, release workflow |
| [release-compatibility.md](.agents/rules/release-compatibility.md) | handshake and public API compatibility |
| [host-access.md](.agents/rules/host-access.md) | capability and authorization model |
| [synchronization.md](.agents/rules/synchronization.md) | upstream host sync checklist |

## Skills Index

| Skill | Use when |
| --- | --- |
| [add-plugin-service](.agents/skills/add-plugin-service/SKILL.md) | exposing a new gRPC service |
| [build-plugin-package](.agents/skills/build-plugin-package/SKILL.md) | building and validating a plugin package |

## Notes Index

| File | Topic |
| --- | --- |
| [2026-10-23-sdk-bootstrap.md](.agents/notes/2026-10-23-sdk-bootstrap.md) | initial SDK scope and layout |
| [2026-10-23-plugin-system.md](.agents/notes/2026-10-23-plugin-system.md) | authorization, extension point, and root layout decisions |
| [2026-10-29-action-catalog.md](.agents/notes/2026-10-29-action-catalog.md) | action.json contract: types, loading, and validation parity with the host |

## Commands

| Command | Purpose |
| --- | --- |
| `make verify` | format check, vet, root tests |
| `make test` | SDK package tests |
| `make test-race` | race tests |
| `make sync-check` | repository-local agent and sync checks |

## Conventions

- Agent docs are English; READMEs and Go comments are Chinese.
- README is Chinese only; do not add `README.zh.md` or a changelog.
- Preserve handshake protocol `1` and cookie values unless modality changes the
  host contract in a coordinated release.
- Do not commit unless explicitly requested.
