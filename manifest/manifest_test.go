// manifest 单测:清单校验、UI 与路由声明。
package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{
		"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin",
		"min_host_version":">=1.0.0 <2.0.0",
		"services":["demo.Greeter"],
		"ui":{"entry":"ui/index.html","routes_file":"ui/routes.json"}
	}`, true)
	got, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if got.Name != "demo-plugin" || got.Services[0] != "demo.Greeter" {
		t.Fatalf("manifest = %+v", got)
	}
}

func TestLoadManifestFieldErrors(t *testing.T) {
	dir := writePlugin(t, "Demo", `{"name":"Demo","version":"bad","binary":"missing","services":[]}`, false)
	_, err := Load(dir)
	var validation *ValidationError
	if !errors.As(err, &validation) || !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("Load() = %v, want ValidationError", err)
	}
	want := map[string]bool{"name": false, "version": false, "binary": false, "services": false}
	for _, fieldErr := range validation.Errors {
		if _, ok := want[fieldErr.Field]; ok {
			want[fieldErr.Field] = true
		}
	}
	for field, found := range want {
		if !found {
			t.Errorf("缺少字段 %s 的错误: %+v", field, validation.Errors)
		}
	}
}

func TestLoadRoutes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "routes.json")
	if err := os.WriteFile(path, []byte(`[{"view":"console","path":"/home","name":"Home","app":"demo-app","action":"demo-list"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	routes, err := LoadRoutes(path)
	if err != nil {
		t.Fatalf("LoadRoutes() = %v", err)
	}
	if len(routes) != 1 || routes[0].Path != "/home" {
		t.Fatalf("routes = %+v", routes)
	}
	if routes[0].View != ViewConsole || routes[0].App != "demo-app" || routes[0].Action != "demo-list" {
		t.Fatalf("扩展字段不符: %+v", routes[0])
	}

	// 契约校验:view 非法、app/action 未成对均拒绝。
	for _, bad := range []string{
		`[{"view":"admin","path":"/a","name":"A"}]`,
		`[{"view":"console","path":"/a","name":"A","app":"demo-app"}]`,
	} {
		if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadRoutes(path); err == nil {
			t.Fatalf("非法路由应被拒绝: %s", bad)
		}
	}
}

func TestDiscover(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"alpha-plugin", "beta-plugin"} {
		target := filepath.Join(root, name)
		writePluginAt(t, target, name, `{"name":"`+name+`","version":"1.0.0","binary":"`+name+`","services":["`+name+`.Service"]}`, true)
	}
	manifests, err := Discover(root)
	if err != nil {
		t.Fatalf("Discover() = %v", err)
	}
	if len(manifests) != 2 || manifests[0].Name != "alpha-plugin" || manifests[1].Name != "beta-plugin" {
		t.Fatalf("manifests = %+v", manifests)
	}
}

func TestCheckHostVersion(t *testing.T) {
	loaded := &Manifest{MinHostVersion: ">=1.2.0 <2.0.0"}
	if err := loaded.CheckHostVersion("1.5.0"); err != nil {
		t.Fatalf("CheckHostVersion(1.5.0) = %v", err)
	}
	if err := loaded.CheckHostVersion("0.9.0"); !errors.Is(err, ErrHostVersionIncompatible) {
		t.Fatalf("CheckHostVersion(0.9.0) = %v, want ErrHostVersionIncompatible", err)
	}
}

func TestLoadManifestGrantsAndExtensions(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{
		"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin",
		"services":["demo.Greeter"],
		"grants":["settings.snapshot.read","upms.user.list"],
		"extensions":[{"name":"greeter","point":"notice.channel","service_name":"demo.Greeter"}]
	}`, true)
	got, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(got.Grants) != 2 {
		t.Fatalf("grants = %+v", got.Grants)
	}
	if len(got.Extensions) != 1 || got.Extensions[0].Point != "notice.channel" {
		t.Fatalf("extensions = %+v", got.Extensions)
	}
}

func TestLoadManifestRejectsUnknownExtensionService(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{
		"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin",
		"services":["demo.Greeter"],
		"extensions":[{"name":"greeter","point":"notice.channel","service_name":"missing.Service"}]
	}`, true)
	_, err := Load(dir)
	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("Load() = %v, want ValidationError", err)
	}
	found := false
	for _, fieldErr := range validation.Errors {
		if fieldErr.Field == "extensions[0].service_name" {
			found = true
		}
	}
	if !found {
		t.Fatalf("errors = %+v", validation.Errors)
	}
}

func writePlugin(t *testing.T, name, data string, complete bool) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	writePluginAt(t, dir, name, data, complete)
	return dir
}

func writePluginAt(t *testing.T, dir, binaryName, data string, complete bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "ui"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, FileName), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if complete {
		if err := os.WriteFile(filepath.Join(dir, binaryName), []byte("bin"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "ui", "index.html"), []byte("ui"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "ui", "routes.json"), []byte(`[]`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
