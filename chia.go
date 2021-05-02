package main

type NetworkInfo struct {
	NetworkName   string `json:"network_name"`
	NetworkPrefix string `json:"network_prefix"`
	Success       bool
}

// Chia node types from server/outbound_message.py
const (
	NodeTypeNone = iota
	NodeTypeFullNode
	NodeTypeHarvester
	NodeTypeFarmer
	NodeTypeTimelord
	NodeTypeIntroducer
	NodeTypeWallet
)

type NodeType int

type Connections struct {
	Connections []struct {
		BytesRead       int     `json:"bytes_read"`
		BytesWritten    int     `json:"bytes_written"`
		CreationTime    float64 `json:"creation_time"`
		LastMessageTime float64 `json:"last_message_time"`
		LocalPort       int     `json:"local_port"`
		NodeId          string  `json:"node_id"`
		PeakHash        string  `json:"peak_hash"`
		PeakHeight      int     `json:"peak_height"`
		PeakWeight      int     `json:"peak_weight"`
		PeerHost        string  `json:"peer_host"`
		PeerPort        int     `json:"peer_port"`
		PeerServerPort  int     `json:"peer_server_port"`
		Type            NodeType
	}
	Success bool
}
