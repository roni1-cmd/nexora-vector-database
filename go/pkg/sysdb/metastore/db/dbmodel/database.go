package dbmodel

import (
	"time"

	"github.com/chroma-core/chroma/go/pkg/types"
)

type Database struct {
	ID        string          `gorm:"id;primaryKey;unique"`
	Name      string          `gorm:"name;type:varchar(128);not null;uniqueIndex:idx_tenantid_name"`
	TenantID  string          `gorm:"tenant_id;type:varchar(128);not null;uniqueIndex:idx_tenantid_name"`
	Ts        types.Timestamp `gorm:"ts;type:bigint;default:0"`
	IsDeleted bool            `gorm:"is_deleted;type:bool;default:false"`
	CreatedAt time.Time       `gorm:"created_at;type:timestamp;not null;default:current_timestamp"`
	UpdatedAt time.Time       `gorm:"updated_at;type:timestamp;not null;default:current_timestamp"`
}

func (v Database) TableName() string {
	return "databases"
}

//go:generate mockery --name=IDatabaseDb
type IDatabaseDb interface {
	GetDatabases(tenantID string, databaseName string) ([]*Database, error)
	ListDatabases(limit *int32, offset *int32, tenantID string) ([]*Database, error)
	Insert(in *Database) error
	DeleteAll() error
	SoftDelete(databaseID string) error
	FinishDatabaseDeletion(cutoffTime time.Time) (uint64, error)
}
