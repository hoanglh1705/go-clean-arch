package repository

import (
	"go-clean-arch/config"
	"go-clean-arch/helper-libs/copyhelper"
	"go-clean-arch/helper-libs/sqlxhelper"
)

type (
	baseRepository struct {
		cfg          *config.Config
		objectCopier copyhelper.ObjectCopier
		sqldb        sqlxhelper.SqlDatabase
	}
)

func NewBaseRepository(
	cfg *config.Config,
	objectCopier copyhelper.ObjectCopier,
	sqldb sqlxhelper.SqlDatabase,
) BaseRepository {
	return &baseRepository{
		cfg:          cfg,
		objectCopier: objectCopier,
		sqldb:        sqldb,
	}
}
