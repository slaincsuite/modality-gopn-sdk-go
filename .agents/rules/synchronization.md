# Upstream Synchronization

The upstream repository is `../modality`. Its AGENTS knowledge and
`library/gopn` contract can change; this SDK must not drift.

## Checklist

When modality changes any of:

1. `AGENTS.md` or `.agents/rules/` guidance useful to plugin or SDK authors;
2. `library/gopn/runner` handshake constants;
3. `library/gopn/manifest` fields or validation;
4. `library/gopn/access` capability semantics;
5. `library/gopn/extensionpoint` registration semantics;
6. plugin management or authorization behavior;

then update this SDK:

- [ ] Filter host-only guidance and import only concepts relevant to plugin authors.
- [ ] Update `core`, `manifest`, `host`, or `extensionpoint` models and tests.
- [ ] Update the Chinese README for public-facing changes.
- [ ] Update `AGENTS.md` / `.agents/` when downstream procedures change.
- [ ] Run `make verify`; add race or subprocess tests for lifecycle changes.

`make sync-check` verifies repository-local agent files and this checklist.
