package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
)

func TestPubSubNotify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hosts := getNetHosts(t, ctx, 2)

	// initialize id services
	{
		ids1, err := identify.NewIDService(hosts[0])
		if err != nil {
			t.Fatal(err)
		}
		ids1.Start()
		defer ids1.Close()

		ids2, err := identify.NewIDService(hosts[1])
		if err != nil {
			t.Fatal(err)
		}
		ids2.Start()
		defer ids2.Close()
	}

	psubs0 := getPubsub(ctx, hosts[0])
	connect(t, hosts[0], hosts[1])
	<-time.After(time.Millisecond * 100)
	psubs1 := getPubsub(ctx, hosts[1])

	// pubsub 0 join topic "test"
	topic0, err := psubs0.Join("test")
	if err != nil {
		t.Fatal(err)
	}
	defer topic0.Close()

	sub0, err := topic0.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	defer sub0.Cancel()

	// pubsub 1 join topic "test"
	topic1, err := psubs1.Join("test")
	if err != nil {
		t.Fatal(err)
	}
	defer topic1.Close()

	sub1, err := topic1.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	defer sub1.Cancel()

	<-time.After(time.Millisecond * 100)

	if len(topic1.ListPeers()) == 0 {
		t.Fatalf("topic1 should at least have 1 peer")
	}

	if len(topic0.ListPeers()) == 0 {
		t.Fatalf("topic0 should at least have 1 peer")
	}
}
