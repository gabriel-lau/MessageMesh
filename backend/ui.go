package backend

import (
	"MessageMesh/debug"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func UIDataLoop(network Network, ctx context.Context) {
	debug.Log("ui", "Wails events emitter started")
	if debug.IsHeadless {

		// Send the user's peer ID once to the frontend and then remove the event listener
		runtime.EventsEmit(ctx, "getUserPeerID", network.P2p.Host.ID())
		runtime.EventsEmit(ctx, "getPeerList", network.ChatRoom.PeerList())
		for {
			select {
			case msg := <-network.ChatRoom.Inbound:
				runtime.EventsEmit(ctx, "getMessage", msg)
				debug.Log("ui", "Message: "+msg.Message)
			case peerIDs := <-network.ChatRoom.PeerIDs:
				runtime.EventsEmit(ctx, "getPeerList", peerIDs)
				debug.Log("ui", "Peers: "+string(len(peerIDs)))
			}
		}
	}
}
