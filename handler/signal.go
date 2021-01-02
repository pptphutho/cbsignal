package handler

import (
	"cbsignal/client"
	"cbsignal/hub"
	"github.com/lexkong/log"
)

type SignalHandler struct {
	Msg   *SignalMsg
	Cli   *client.Client
}

func (s *SignalHandler)Handle() {
	h := hub.GetInstance()
	//log.Infof("load client Msg %v", s.Msg)
	target, ok := h.Clients.Load(s.Msg.ToPeerId) //判断节点是否还在线
	if ok {
		//log.Infof("found client %s", s.Msg.ToPeerId)
		resp := SignalResp{
			Action: "signal",
			FromPeerId: s.Cli.PeerId,
			Data: s.Msg.Data,
		}
		hub.SendJsonToClient(s.Msg.ToPeerId, resp, true)
		if !target.(*client.Client).LocalNode {
			log.Warnf("send signal msg from %s to %s on node %s", s.Cli.PeerId, s.Msg.ToPeerId, target.(*client.Client).RpcNodeAddr)
		}
	} else {
		log.Infof("Peer %s not found, ", s.Msg.ToPeerId)
		resp := SignalResp{
			Action: "signal",
			FromPeerId: s.Msg.ToPeerId,
		}
		// 发送一次后，同一peerId下次不再发送，节省带宽
		if !s.Cli.InvalidPeers[s.Msg.ToPeerId] {
			s.Cli.InvalidPeers[s.Msg.ToPeerId] = true
			hub.SendJsonToClient(s.Cli.PeerId, resp, true)
		}
	}
}


