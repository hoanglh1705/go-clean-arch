package sqlormhelper

import (
	"context"
	"fmt"

	"go-clean-arch/helper-libs/loghelper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormMysqlDB struct {
	db *gorm.DB
}

func NewGormMysqlDB(options *GormConnectionOptions) SqlGormDatabase {
	db, err := initGormMysqlDB(options)
	if err != nil {
		loghelper.Logger.Panic("Failed to init mysql", err)
	}
	return &gormMysqlDB{
		db: db,
	}
}

func (h gormMysqlDB) Open() *gorm.DB {
	return h.db
}

func (h *gormMysqlDB) Close() error {
	return nil
}

func (h gormMysqlDB) Begin() *gorm.DB {
	return h.db.Begin()
}

func (h gormMysqlDB) Commit(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (h gormMysqlDB) Rollback(tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}

func (h gormMysqlDB) GetConn() (*gorm.DB, error) {
	return h.db, nil
}

func (h gormMysqlDB) List(ctx context.Context, input *ListParams) *gorm.DB {
	db := h.Open()
	result := db.Model(input.Model)

	if input.LastQueryId != "" {
		if input.Descending {
			result.Where("id < ?", input.LastQueryId)
		} else {
			result.Where("id > ?", input.LastQueryId)
		}
	}
	if len(input.WhereClauses) > 0 {
		for _, clause := range input.WhereClauses {
			result.Where(clause)
		}
	}
	if input.Descending {
		if input.OrderByColumn != "" {
			result.Order(clause.OrderByColumn{Column: clause.Column{Name: input.OrderByColumn}, Desc: true})
		} else {
			result.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true})
		}
	}
	if input.Limit > 0 {
		result.Limit(input.Limit)
	}
	return result.Find(input.Target)
}

func initGormMysqlDB(options *GormConnectionOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), options.GormConfig)
	if err != nil {
		return nil, err
	}
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if options.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(options.MaxOpenConns)
	}

	loghelper.Logger.Infof("Gorm Mysql: Successfully connected!")
	return db, nil
}
