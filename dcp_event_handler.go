package dcpsql

import (
	"github.com/Trendyol/go-dcp-sql/sql/bulk"
)

type DcpEventHandler struct {
	bulk     *bulk.Bulk
	isFinite bool
}

func (h *DcpEventHandler) BeforeRebalanceStart() {
}

func (h *DcpEventHandler) AfterRebalanceStart() {
}

func (h *DcpEventHandler) BeforeRebalanceEnd() {
}

func (h *DcpEventHandler) AfterRebalanceEnd() {
}

func (h *DcpEventHandler) BeforeStreamStart() {
	if h.isFinite {
		return
	}
	h.bulk.PrepareEndRebalancing()
}

func (h *DcpEventHandler) AfterStreamStart() {
}

func (h *DcpEventHandler) BeforeStreamStop() {
	h.bulk.PrepareStartRebalancing()
}

func (h *DcpEventHandler) AfterStreamStop() {
}
