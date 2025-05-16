package server

import (
	"github.com/glin-gogogo/boxo/bitswap/server/internal/decision"
)

type (
	Receipt                = decision.Receipt
	PeerBlockRequestFilter = decision.PeerBlockRequestFilter
	TaskComparator         = decision.TaskComparator
	TaskInfo               = decision.TaskInfo
	ScoreLedger            = decision.ScoreLedger
	ScorePeerFunc          = decision.ScorePeerFunc
	PeerLedger             = decision.PeerLedger
	PeerEntry              = decision.PeerEntry
)
