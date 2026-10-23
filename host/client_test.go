// client_test.go 覆盖宿主能力授权与默认拒绝行为。
package host

import (
	"context"
	"errors"
	"testing"
)

type fakeCaller struct {
	authorization Authorization
	request       Request
}

func (f *fakeCaller) Call(_ context.Context, a Authorization, r Request) (Response, error) {
	f.authorization = a
	f.request = r
	return Response{Data: "ok"}, nil
}

func TestClientAuthorizeAndCall(t *testing.T) {
	capability := Capability{Name: "upms.user.list", Module: "upms", Mode: AccessRead}
	caller := &fakeCaller{}
	client := New(Config{PluginName: "demo-plugin", Capabilities: []Capability{capability}}, caller)
	got, err := client.Call(context.Background(), capability.Name, map[string]int64{"id": 7})
	if err != nil {
		t.Fatalf("Call() = %v", err)
	}
	if got.Data != "ok" || caller.request.Capability != capability.Name {
		t.Fatalf("call result = %+v, request = %+v", got, caller.request)
	}
	if caller.authorization.Plugin != "demo-plugin" {
		t.Fatalf("authorization = %+v", caller.authorization)
	}
}

func TestClientDefaultDeny(t *testing.T) {
	client := New(Config{PluginName: "demo-plugin"}, &fakeCaller{})
	if _, err := client.Call(context.Background(), "plugin.core.write", nil); !errors.Is(err, ErrForbidden) {
		t.Fatalf("Call() = %v, want ErrForbidden", err)
	}
	if _, err := client.Call(context.Background(), "missing", nil); !errors.Is(err, ErrForbidden) {
		t.Fatalf("Call() = %v, want ErrForbidden", err)
	}
}

func TestClientRejectsInvalidCapability(t *testing.T) {
	client := New(Config{PluginName: "demo-plugin", Capabilities: []Capability{{Name: "bad"}}}, &fakeCaller{})
	if _, err := client.Call(context.Background(), "bad", nil); !errors.Is(err, ErrInvalidCapability) {
		t.Fatalf("Call() = %v, want ErrInvalidCapability", err)
	}
}
