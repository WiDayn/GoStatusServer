package model

type PingRecords struct {
	CT PingRecord
	CU PingRecord
	CM PingRecord
}

type PingRecord struct {
	PacketsSent    int
	PacketsReceive int
	Time           []string
	AvgRTT         []int64
}

type BasicRecords struct {
	CPUAvg                []float64
	MemUsedPercent        []float64
	DiskPercent           []uint64
	NowUpStreamDataSize   []int
	NowDownStreamDataSize []int
	Time                  []string
}
