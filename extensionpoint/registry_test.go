// registry_test.go 覆盖扩展声明注册与重复拒绝。
package extensionpoint

import "testing"

func TestRegistryRegisterAndList(t *testing.T) {
	registry := NewRegistry()
	extension := Extension{Name: "custom", Point: "notice.channel", ServiceName: "starter.v1.Service"}
	if err := registry.Register(extension); err != nil {
		t.Fatalf("Register() = %v", err)
	}
	if got := registry.Extensions(); len(got) != 1 || got[0] != extension {
		t.Fatalf("Extensions() = %+v", got)
	}
}

func TestRegistryRejectsInvalidAndDuplicate(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(Extension{Name: "bad"}); err == nil {
		t.Fatal("invalid extension should be rejected")
	}
	extension := Extension{Name: "custom", Point: "notice.channel", ServiceName: "starter.v1.Service"}
	if err := registry.Register(extension); err != nil {
		t.Fatalf("Register() = %v", err)
	}
	if err := registry.Register(extension); err == nil {
		t.Fatal("duplicate extension should be rejected")
	}
}
