---
name: build-plugin-package
description: Build and validate a distributable modality plugin directory.
---

1. Build the binary for the target host platform:
   `go build -trimpath -o <dir>/<name> ./cmd/<name>`.
2. Place `plugin.json` beside the binary. `name` must be kebab-case and equal
   the directory name; `binary` must equal the compiled file name.
3. Ensure `services` exactly matches the registered `modality.Service.Name`
   values and is non-empty and unique.
4. Add `ui/index.html` only when `ui` is declared; add `routes.json` only when
   it is valid JSON with `path` and `name` on every entry.
5. Declare applications, actions, and resources in `action.json` beside
   `plugin.json` (the `actions` field may point elsewhere). Validate it with
   `manifest.LoadActions(path)` when it exists: application codes and action
   names are kebab-case and unique, `kind` is one of `read|write|list|other`,
   and every `resource` references a declared `resources[].code`. The same
   file is the plugin frontend's permission source; never copy the action
   list into the frontend.
6. Validate with `manifest.Load(dir)` (it loads `action.json` too) and
   `Manifest.CheckHostVersion(hostVersion)`; inspect `ValidationError`
   field paths on failure.
7. Run the plugin package through the sample subprocess test pattern or a host
   launch test before distribution.
8. Do not commit compiled binaries; ship them as release assets.
