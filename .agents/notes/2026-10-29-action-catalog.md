# 2026-10-29 - Action Catalog (action.json)

## Context

The host moved plugin application declarations out of `plugin.json` into a
dedicated `action.json` at the plugin root (host PRD
`prd/modules/plugin/apps.md` FR-1/FR-4/FR-5). The same file is the single
source of truth for the host policy system (sam functions
`应用编码.操作名`) and for the plugin frontend, which now takes its
permissions only from `GET /open-api/v1/upms/policies/my-actions?app=<code>`.
Plugin authors must be able to validate the file in Go before packaging.

## Decisions

1. **Contract parity.** `manifest/actions.go` mirrors the host
   `library/gopn/manifest/actions.go`: exported names `App`, `Action`,
   `Resource`, `ActionFileName`, `OpKindRead`/`OpKindWrite`/`OpKindList`/
   `OpKindOther`, `ValidOperationKind`, and `LoadActions(path)`. The SDK never
   renames the host contract; authors can port code between the two.
2. **Storage and discovery.** `plugin.json` gains an optional `actions`
   pointer (relative path); when omitted, a root `action.json` is discovered.
   Neither present means a service-only plugin, and `Manifest.Apps` stays
   empty. `Load` fills `Manifest.Apps` from the catalog; a declared pointer
   whose file is missing fails validation with field `actions`.
3. **Field-path validation.** Errors aggregate in the existing
   `ValidationError`/`FieldError` shape with paths such as
   `apps[0].actions[2].kind`, matching the host contract and the package's
   established style.
4. **Sentinel error.** `ErrInvalidActions` marks catalog read/parse/validation
   failures; `Load` additionally wraps them with `ErrInvalidManifest` so both
   `errors.Is` checks hold for an invalid plugin package.
5. **Kinds are permission categories.** `read`/`write`/`list`/`other`
   (读/写/列表/其他) are the plugin-defined permission categories; the
   frontend convention is 读=打开/查看, 写=编辑, 列表=看列表/导出,
   其他=其他入口. The SDK only models and validates them; decision data comes
   from the host endpoint.

## Current State

- `manifest` covers `plugin.json`, `routes.json`, and `action.json`;
  `make verify` passes with table-driven tests for valid catalogs, field-path
  errors, discovery, declared paths, and invalid catalogs.
- README (root and `manifest/`) documents the file; the build-plugin-package
  skill includes the validation step.

## Open Items

- Publish the next SDK version before the sample starter can validate its
  `action.json` through `manifest.LoadActions` in `make verify`.