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
	// Info returns address information about the host.
	Info func(context.Context) (peer.AddrInfo, error) `perm:"admin"`
	// Peers returns connected peers.
	Peers func(context.Context) ([]peer.ID, error) `perm:"admin"`
	// PeerInfo returns a small slice of information Peerstore has on the
	// given peer.
	PeerInfo func(ctx context.Context, id peer.ID) (peer.AddrInfo, error) `perm:"admin"`
	// Connect ensures there is a connection between this host and the peer with
	// given peer.
	Connect func(ctx context.Context, pi peer.AddrInfo) error `perm:"admin"`
	// ClosePeer closes the connection to a given peer.
	ClosePeer func(ctx context.Context, id peer.ID) error `perm:"admin"`
	// Connectedness returns a state signaling connection capabilities.
	Connectedness func(ctx context.Context, id peer.ID) (network.Connectedness, error) `perm:"admin"`
	// NATStatus returns the current NAT status.
	NATStatus func(context.Context) (network.Reachability, error) `perm:"admin"`
	// BlockPeer adds a peer to the set of blocked peers.
	BlockPeer func(ctx context.Context, p peer.ID) error `perm:"admin"`
	// UnblockPeer removes a peer from the set of blocked peers.
	UnblockPeer func(ctx context.Context, p peer.ID) error `perm:"admin"`
	// ListBlockedPeers returns a list of blocked peers.
	ListBlockedPeers func(context.Context) ([]peer.ID, error) `perm:"admin"`
	// Protect adds a peer to the list of peers who have a bidirectional
	// peering agreement that they are protected from being trimmed, dropped
	// or negatively scored.
	Protect func(ctx context.Context, id peer.ID, tag string) error `perm:"admin"`
	// Unprotect removes a peer from the list of peers who have a bidirectional
	// peering agreement that they are protected from being trimmed, dropped
	// or negatively scored, returning a bool representing whether the given
	// peer is protected or not.
	Unprotect func(ctx context.Context, id peer.ID, tag string) (bool, error) `perm:"admin"`
	// IsProtected returns whether the given peer is protected.
	IsProtected func(ctx context.Context, id peer.ID, tag string) (bool, error) `perm:"admin"`
	// BandwidthStats returns a Stats struct with bandwidth metrics for all
	// data sent/received by the local peer, regardless of protocol or remote
	// peer IDs.
	BandwidthStats func(context.Context) (metrics.Stats, error) `perm:"admin"`
	// BandwidthForPeer returns a Stats struct with bandwidth metrics associated with the given peer.ID.
	// The metrics returned include all traffic sent / received for the peer, regardless of protocol.
	BandwidthForPeer func(ctx context.Context, id peer.ID) (metrics.Stats, error) `perm:"admin"`
	// BandwidthForProtocol returns a Stats struct with bandwidth metrics associated with the given
	// protocol.ID.
	BandwidthForProtocol func(ctx context.Context, proto protocol.ID) (metrics.Stats, error) `perm:"admin"`
	// ResourceState returns the state of the resource manager.
	ResourceState func(context.Context) (rcmgr.ResourceManagerStat, error) `perm:"admin"`
	// PubSubPeers returns the peer IDs of the peers joined on
	// the given topic.
	PubSubPeers func(ctx context.Context, topic string) ([]peer.ID, error) `perm:"admin"`
}
