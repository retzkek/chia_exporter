package main

type NetworkInfo struct {
	NetworkName   string `json:"network_name"`
	NetworkPrefix string `json:"network_prefix"`
	Success       bool
}

//TODO figure out what type some of these fields should be, set to interface{} for now
type BlockchainState struct {
	BlockchainState struct {
		Difficulty                  int
		GenesisChallengeInitialized bool `json:"genesis_challenge_initialized"`
		MempoolSize                 int  `json:"mempool_size"`
		Peak                        struct {
			ChallengeBlockInfoHash string `json:"challenge_block_info_hash"`
			ChallengeVDFOutput     struct {
				Data string
			} `json:"challenge_vdf_output"`
			Deficit                            int
			FarmerPuzzleHash                   string `json:"farmer_puzzle_hash"`
			Fees                               float64
			FinishedChallengeSlotHashes        interface{} `json:"finished_challenge_slot_hashes"`
			FinishedInfusedChallengeSlotHashes interface{} `json:"finished_infused_challenge_slot_hashes"`
			HeaderHash                         string      `json:"header_hash"`
			Height                             int
			InfusedChallengeVDFOutput          struct {
				Date string
			} `json:"infused_challenge_vdf_output"`
			Overflow                   bool
			PoolPuzzleHash             string      `json:"pool_puzzle_hash"`
			PrevHash                   string      `json:"prev_hash"`
			PrevTransactionBlockHash   string      `json:"prev_transaction_block_hash"`
			PrevTransactionBlockHeight int         `json:"prev_transaction_block_height"`
			RequiredIters              int         `json:"required_iters"`
			RewardClaimsIncorporated   interface{} `json:"reward_claims_incorporated"`
			RewardInfusionNewChallenge string      `json:"reward_infusion_new_challenge"`
			SignagePointIndex          int         `json:"signage_point_index"`
			SubEpochSummaryIncluded    interface{} `json:"sub_epoch_summary_included"`
			SubSlotIters               int         `json:"sub_slot_iters"`
			Timestamp                  interface{}
			TotalIters                 int `json:"total_iters"`
			Weight                     int
		}
		Space        float64
		SubSlotIters int `json:"sub_slot_iters"`
		Sync         struct {
			SyncMode           bool `json:"sync_mode"`
			SyncProgressHeight int  `json:"sync_progress_height"`
			SyncTipHeight      int  `json:"sync_tip_height"`
			Synced             bool
		}
	} `json:"blockchain_state"`
	Success bool
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
	NumNodeTypes = 6
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

type Wallet struct {
	ID        int
	Name      string
	Type      int
	Data      string
	StringID  string
	PublicKey string
}

type Wallets struct {
	Wallets []Wallet
	Success bool
}

type WalletBalance struct {
	WalletBalance struct {
		ConfirmedBalance   int64 `json:"confirmed_wallet_balance"`
		MaxSendAmount      int64 `json:"max_send_amount"`
		PendingChange      int64 `json:"pending_change"`
		SpendableBalance   int64 `json:"spendable_balance"`
		UnconfirmedBalance int64 `json:"unconfirmed_wallet_balance"`
		WalletID           int   `json:"wallet_id"`
	} `json:"wallet_balance"`
	Success bool
}

type WalletSyncStatus struct {
	GenesisInitialized bool `json:"genesis_initialized"`
	Synced             bool
	Syncing            bool
	Succes             bool
}

type WalletHeightInfo struct {
	Height  int64
	Success bool
}

type WalletPublicKeys struct {
	PublicKeyFingerprints []int `json:"public_key_fingerprints"`
	Success               bool
}

type FarmedAmount struct {
	FarmedAmount     int64 `json:"farmed_amount"`
	RewardAmount     int64 `json:"farmer_reward_amount"`
	FeeAmount        int64 `json:"fee_amount"`
	LastHeightFarmed int64 `json:"last_height_farmed"`
	PoolRewardAmount int64 `json:"pool_reward_amount"`
	Success          bool
}

type PoolState struct {
	PoolState []struct {
		CurrentDificulty      int64        `json:"current_difficulty"`
		CurrentPoints         int64        `json:"current_points"`
		PointsAcknowledged24h [][2]float64 `json:"points_acknowledged_24h"`
		PointsFound24h        [][2]float64 `json:"points_found_24h"`
		PoolConfig            struct {
			LauncherId string `json:"launcher_id"`
			PoolURL    string `json:"pool_url"`
		} `json:"pool_config"`
	} `json:"pool_state"`
	Success bool
}

type PlotData struct {
	FileSize      int64   `json:"file_size"`
	Filename      string  `json:"filename"`
	PlotSeed      string  `json:"plot-seed"`
	PlotID        string  `json:"plot_id"`
	PublicKey     string  `json:"plot_public_key"`
	PoolContract  string  `json:"pool_contract_puzzle_hash"`
	PoolPublicKey string  `json:"pool_public_key"`
	Size          int64   `json:"size"`
	TimeModified  float64 `json:"time_modified"`
}

type PlotFiles struct {
	FailedToOpen []string   `json:"failed_to_open_filenames"`
	NotFound     []string   `json:"not_found_filenames"`
	Plots        []PlotData `json:"plots"`
	Success      bool
}
