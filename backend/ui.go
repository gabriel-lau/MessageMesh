package backend

import (
	"MessageMesh/debug"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func UIDataLoop(network Network, ctx context.Context) {
	debug.Log("ui", "Wails events emitter started")
	if debug.GetEnvVar("HEADLESS") == "false" {
		for {
			select {
			case msg := <-network.ChatRoom.Inbound:
				runtime.EventsEmit(ctx, "getMessage", msg.Message)
				debug.Log("ui", "Message: "+msg.Message)
			case peerIDs := <-network.ChatRoom.PeerIDs:
				runtime.EventsEmit(ctx, "getPeersList", len(peerIDs))
				debug.Log("ui", "Peers: "+string(len(peerIDs)))
			}
		}
	}
}
