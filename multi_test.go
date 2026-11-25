package multiresolver

import (
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

func init() {
	resolver.Register(startupBuilder{})
}

type startupBuilder struct{}

func (startupBuilder) Scheme() string { return "startup-test" }

func (b startupBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	group := cc.(*partialClientConn).parent
	if got := len(group.parts); got != 2 {
		return nil, fmt.Errorf("child can publish with %d of 2 child connections ready", got)
	}
	child := manual.NewBuilderWithScheme(b.Scheme())
	child.InitialState(resolver.State{Addresses: []resolver.Address{{Addr: target.Endpoint()}}})
	return child.Build(target, cc, opts)
}

type startupClientConn struct {
	resolver.ClientConn
	state resolver.State
}

func (cc *startupClientConn) UpdateState(s resolver.State) error {
	cc.state = s
	return nil
}

func TestBuildPreparesChildrenBeforeFirstUpdate(t *testing.T) {
	cc := &startupClientConn{}
	target := ParseTarget("multi:///startup-test:///192.0.2.1:443,startup-test:///192.0.2.2:443")
	r, err := (builder{}).Build(target, cc, resolver.BuildOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	addresses := []resolver.Address{{Addr: "192.0.2.1:443"}, {Addr: "192.0.2.2:443"}}
	if !reflect.DeepEqual(cc.state.Addresses, addresses) {
		t.Errorf("published addresses = %v, want %v", cc.state.Addresses, addresses)
	}
}
