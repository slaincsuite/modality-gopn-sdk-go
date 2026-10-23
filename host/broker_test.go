// broker_test.go 覆盖宿主 broker 通道的空值防护。
package host

import (
	"context"
	"errors"
	"testing"
)

func TestBrokerCallerRequiresHost(t *testing.T) {
	if _, err := NewBrokerCaller(nil); !errors.Is(err, ErrHostUnavailable) {
		t.Fatalf("NewBrokerCaller(nil) = %v", err)
	}
	var caller *BrokerCaller
	if _, err := caller.Call(context.Background(), Authorization{}, Request{}); !errors.Is(err, ErrHostUnavailable) {
		t.Fatalf("Call() = %v", err)
	}
	if err := caller.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}
}
