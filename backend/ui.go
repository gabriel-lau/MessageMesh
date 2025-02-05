package backend

import (
	"MessageMesh/debug"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func UIDataLoop(network Network, ctx context.Context) {
	debug.Log("ui", "Wails events emitter started")
	if !debug.IsHeadless {

		// Send the user's peer ID once to the frontend and then remove the event listener
		runtime.EventsEmit(ctx, "getUserPeerID", network.P2pService.Host.ID())
		runtime.EventsEmit(ctx, "getPeerList", network.PubSubService.PeerList())
		for {
			select {
			case msg := <-network.PubSubService.Inbound:
				runtime.EventsEmit(ctx, "getMessage", msg)
				debug.Log("ui", "Message: "+msg.Message)
			case peerIDs := <-network.PubSubService.PeerIDs:
				runtime.EventsEmit(ctx, "getPeerList", peerIDs)
				debug.Log("ui", "Peers: "+string(len(peerIDs)))
			}
		}
	}
}
