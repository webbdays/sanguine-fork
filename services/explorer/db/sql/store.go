package sql

import (
	"context"
	"fmt"
	"github.com/synapsecns/sanguine/core/dbcommon"
	gormClickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"time"
)

// Store is the clickhouse store. It extends the base store for sqlite specific queries.
type Store struct {
	db *gorm.DB
}

// UNSAFE_DB gets the underlying gorm db.
//
//nolint:golint,revive,stylecheck
func (s *Store) UNSAFE_DB() *gorm.DB {
	return s.db
}

// OpenGormClickhouse opens a gorm connection to clickhouse.
func OpenGormClickhouse(ctx context.Context, address string) (*Store, error) {
	clickhouseDB, err := gorm.Open(gormClickhouse.New(gormClickhouse.Config{
		DSN: address,
	}), &gorm.Config{
		Logger:               dbcommon.GetGormLogger(logger),
		FullSaveAssociations: true,
		NowFunc:              time.Now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm clickhouse: %w", err)
	}

	// load all models
	err = clickhouseDB.WithContext(ctx).Set("gorm:table_options", "ENGINE=ReplacingMergeTree(insert_time) ORDER BY (event_index, block_number, event_type, tx_hash, chain_id, contract_address)").AutoMigrate(&SwapEvent{}, &BridgeEvent{})
	if err != nil {
		return nil, fmt.Errorf("could not migrate on clickhouse: %w", err)
	}
	return &Store{clickhouseDB}, nil
}