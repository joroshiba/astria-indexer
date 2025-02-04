// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Action -
type Action struct {
	*postgres.Table[*storage.Action]
}

// NewAction -
func NewAction(db *database.Bun) *Action {
	return &Action{
		Table: postgres.NewTable[*storage.Action](db),
	}
}

func (a *Action) ByBlock(ctx context.Context, height types.Level, limit, offset int) (actions []storage.ActionWithTx, err error) {
	query := a.DB().NewSelect().
		Model((*storage.Action)(nil)).
		Where("height = ?", height)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = a.DB().NewSelect().
		TableExpr("(?) as action", query).
		ColumnExpr("action.*").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = action.tx_id").
		Scan(ctx, &actions)
	return
}

func (a *Action) ByTxId(ctx context.Context, txId uint64, limit, offset int) (actions []storage.Action, err error) {
	query := a.DB().NewSelect().
		Model(&actions).
		Where("tx_id = ?", txId)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)
	return
}

func (a *Action) ByAddress(ctx context.Context, addressId uint64, filters storage.AddressActionsFilter) (actions []storage.AddressAction, err error) {
	query := a.DB().NewSelect().
		Model((*storage.AddressAction)(nil)).
		Where("address_id = ?", addressId)

	if filters.ActionTypes.Bits > 0 {
		query = query.Where("action_type IN (?)", bun.In(filters.ActionTypes.Strings()))
	}

	query = sortScope(query, "action_id", filters.Sort)
	query = limitScope(query, filters.Limit)
	query = offsetScope(query, filters.Offset)

	err = a.DB().NewSelect().
		TableExpr("(?) as address_action", query).
		ColumnExpr("address_action.*").
		ColumnExpr("action.id as action__id, action.height as action__height, action.time as action__time, action.position as action__position, action.type as action__type, action.tx_id as action__tx_id, action.data as action__data").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = address_action.tx_id").
		Join("left join action on action.id = address_action.action_id").
		Scan(ctx, &actions)
	return
}

func (a *Action) ByRollup(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) (actions []storage.RollupAction, err error) {
	query := a.DB().NewSelect().
		Model((*storage.RollupAction)(nil)).
		Where("rollup_id = ?", rollupId)

	query = sortScope(query, "action_id", sort)
	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = a.DB().NewSelect().
		TableExpr("(?) as rollup_action", query).
		ColumnExpr("rollup_action.*").
		ColumnExpr("action.id as action__id, action.height as action__height, action.time as action__time, action.position as action__position, action.type as action__type, action.tx_id as action__tx_id, action.data as action__data").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = rollup_action.tx_id").
		Join("left join action on action.id = rollup_action.action_id").
		Scan(ctx, &actions)
	return
}
