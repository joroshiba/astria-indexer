// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/pkg/errors"
)

type Stats struct {
	db *database.Bun
}

func NewStats(conn *database.Bun) Stats {
	return Stats{conn}
}

func (s Stats) Series(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewBlockStatsByHour
	case storage.TimeframeDay:
		view = storage.ViewBlockStatsByDay
	case storage.TimeframeMonth:
		view = storage.ViewBlockStatsByMonth
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view)

	switch name {
	case storage.SeriesDataSize:
		query.ColumnExpr("ts, data_size as value")
	case storage.SeriesTPS:
		query.ColumnExpr("ts, tps as value, tps_max as max, tps_min as min")
	case storage.SeriesBPS:
		query.ColumnExpr("ts, bps as value, bps_max as max, bps_min as min")
	case storage.SeriesRBPS:
		query.ColumnExpr("ts, rbps as value, rbps_max as max, rbps_min as min")
	case storage.SeriesFee:
		query.ColumnExpr("ts, fee as value")
	case storage.SeriesSupplyChange:
		query.ColumnExpr("ts, supply_change as value")
	case storage.SeriesBlockTime:
		query.ColumnExpr("ts, block_time as value")
	case storage.SeriesTxCount:
		query.ColumnExpr("ts, tx_count as value")
	case storage.SeriesBytesInBlock:
		query.ColumnExpr("ts, bytes_in_block as value")
	case storage.SeriesGasPrice:
		query.ColumnExpr("ts, gas_price as value")
	case storage.SeriesGasEfficiency:
		query.ColumnExpr("ts, gas_efficiency as value")
	case storage.SeriesGasWanted:
		query.ColumnExpr("ts, gas_wanted as value")
	case storage.SeriesGasUsed:
		query.ColumnExpr("ts, gas_used as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("ts < ?", req.To)
	}

	err = query.Limit(100).Scan(ctx, &response)
	return
}

func (s Stats) RollupSeries(ctx context.Context, rollupId uint64, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewRollupStatsByHour
	case storage.TimeframeDay:
		view = storage.ViewRollupStatsByDay
	case storage.TimeframeMonth:
		view = storage.ViewRollupStatsByMonth
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view).
		Where("rollup_id = ?", rollupId)

	switch name {
	case storage.RollupSeriesActionsCount:
		query.ColumnExpr("ts, actions_count as value")
	case storage.RollupSeriesAvgSize:
		query.ColumnExpr("ts, avg_size as value")
	case storage.RollupSeriesMaxSize:
		query.ColumnExpr("ts, max_size as value")
	case storage.RollupSeriesMinSize:
		query.ColumnExpr("ts, min_size as value")
	case storage.RollupSeriesSize:
		query.ColumnExpr("ts, size as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("ts < ?", req.To)
	}

	err = query.Limit(100).Scan(ctx, &response)
	return
}

func (s Stats) Summary(ctx context.Context) (summary storage.NetworkSummary, err error) {
	err = s.db.DB().NewSelect().Table(storage.ViewBlockStatsByMonth).
		ColumnExpr("sum(data_size) as data_size, sum(fee) as fee, sum(supply_change) as supply, sum(tx_count) as tx_count, sum(bytes_in_block) as bytes_in_block").
		ColumnExpr("avg(tps) as tps, avg(bps) as bps, avg(rbps) as rbps, avg(block_time) as block_time").
		Scan(ctx, &summary)
	return
}
