package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
)

type API struct {
	Peers                func(context.Context) ([]peer.ID, error)                             `perm:"admin"`
	PeerInfo             func(ctx context.Context, id peer.ID) (peer.AddrInfo, error)         `perm:"admin"`
	Connect              func(ctx context.Context, pi peer.AddrInfo) error                    `perm:"admin"`
	ClosePeer            func(ctx context.Context, id peer.ID) error                          `perm:"admin"`
	Connectedness        func(ctx context.Context, id peer.ID) (network.Connectedness, error) `perm:"admin"`
	NATStatus            func(context.Context) (network.Reachability, error)                  `perm:"admin"`
	BlockPeer            func(ctx context.Context, p peer.ID) error                           `perm:"admin"`
	UnblockPeer          func(ctx context.Context, p peer.ID) error                           `perm:"admin"`
	ListBlockedPeers     func(context.Context) ([]peer.ID, error)                             `perm:"admin"`
	Protect              func(ctx context.Context, id peer.ID, tag string) error              `perm:"admin"`
	Unprotect            func(ctx context.Context, id peer.ID, tag string) (bool, error)      `perm:"admin"`
	IsProtected          func(ctx context.Context, id peer.ID, tag string) (bool, error)      `perm:"admin"`
	BandwidthStats       func(context.Context) (metrics.Stats, error)                         `perm:"admin"`
	BandwidthForPeer     func(ctx context.Context, id peer.ID) (metrics.Stats, error)         `perm:"admin"`
	BandwidthForProtocol func(ctx context.Context, proto protocol.ID) (metrics.Stats, error)  `perm:"admin"`
	ResourceState        func(context.Context) (rcmgr.ResourceManagerStat, error)             `perm:"admin"`
	PubSubPeers          func(ctx context.Context, topic string) ([]peer.ID, error)           `perm:"admin"`
}
