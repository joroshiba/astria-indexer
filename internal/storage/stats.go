// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Timeframe string

const (
	TimeframeHour  Timeframe = "hour"
	TimeframeDay   Timeframe = "day"
	TimeframeMonth Timeframe = "month"
)

type TPS struct {
	Low               float64
	High              float64
	Current           float64
	ChangeLastHourPct float64
}

type TxCountForLast24hItem struct {
	Time    time.Time `bun:"ts"`
	TxCount int64     `bun:"tx_count"`
	TPS     float64   `bun:"tps"`
}

type SeriesRequest struct {
	From time.Time
	To   time.Time
}

func NewSeriesRequest(from, to int64) (sr SeriesRequest) {
	if from > 0 {
		sr.From = time.Unix(from, 0).UTC()
	}
	if to > 0 {
		sr.To = time.Unix(to, 0).UTC()
	}
	return
}

type SeriesItem struct {
	Time  time.Time `bun:"ts"`
	Value string    `bun:"value"`
	Max   string    `bun:"max"`
	Min   string    `bun:"min"`
}

const (
	SeriesDataSize      = "data_size"
	SeriesTPS           = "tps"
	SeriesBPS           = "bps"
	SeriesRBPS          = "rbps"
	SeriesFee           = "fee"
	SeriesSupplyChange  = "supply_change"
	SeriesBlockTime     = "block_time"
	SeriesTxCount       = "tx_count"
	SeriesBytesInBlock  = "bytes_in_block"
	SeriesGasPrice      = "gas_price"
	SeriesGasUsed       = "gas_used"
	SeriesGasWanted     = "gas_wanted"
	SeriesGasEfficiency = "gas_efficiency"

	RollupSeriesActionsCount = "actions_count"
	RollupSeriesSize         = "size"
	RollupSeriesAvgSize      = "avg_size"
	RollupSeriesMinSize      = "min_size"
	RollupSeriesMaxSize      = "max_size"
)

type NetworkSummary struct {
	DataSize     int64           `bun:"data_size"`
	TPS          float64         `bun:"tps"`
	BPS          float64         `bun:"bps"`
	RBPS         float64         `bun:"rbps"`
	Fee          decimal.Decimal `bun:"fee"`
	Supply       decimal.Decimal `bun:"supply"`
	BlockTime    float64         `bun:"block_time"`
	TxCount      int64           `bun:"tx_count"`
	BytesInBlock int64           `bun:"bytes_in_block"`
}

type RollupSummary struct {
	ActionsCount int64 `bun:"actions_count"`
	Size         int64 `bun:"size"`
	AvgSize      int64 `bun:"avg_size"`
	MinSize      int64 `bun:"min_size"`
	MaxSize      int64 `bun:"max_size"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IStats interface {
	Summary(ctx context.Context) (NetworkSummary, error)
	Series(ctx context.Context, timeframe Timeframe, name string, req SeriesRequest) ([]SeriesItem, error)
	RollupSeries(ctx context.Context, rollupId uint64, timeframe Timeframe, name string, req SeriesRequest) ([]SeriesItem, error)
}
