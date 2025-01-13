package sqlormhelper

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"go-clean-arch/helper-libs/loghelper"

	"github.com/imdatngo/gowhere"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	"gorm.io/gorm/schema"
)

type gormPostgresqlDB struct {
	db *gorm.DB
}

func NewGormPostgresqlDB(options *GormConnectionOptions) SqlGormDatabase {
	pretreatmentConfig(true)
	db, err := initGormPostgresqlDB(options)
	if err != nil {
		loghelper.Logger.Panic("Failed to init postgresql", zap.Error(err))
	}
	return &gormPostgresqlDB{
		db: db,
	}
}

func (h gormPostgresqlDB) Open() *gorm.DB {
	return h.db
}

func (h *gormPostgresqlDB) Close() error {
	return nil
}

func (h gormPostgresqlDB) Begin() *gorm.DB {
	return h.db.Begin()
}

func (h gormPostgresqlDB) Commit(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (h gormPostgresqlDB) Rollback(tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}

func (h gormPostgresqlDB) GetConn() (*gorm.DB, error) {
	return h.db, nil
}

func (h gormPostgresqlDB) CreateConditions(ctx context.Context, expressions []*ConditionExpression) []any {
	conds := []any{}

	if len(expressions) > 0 {
		for _, exp := range expressions {
			switch exp.Operator {
			case "Eq":
				conds = append(conds, clause.Eq{
					Column: exp.Column,
					Value:  exp.Value,
				})
			case "Gt":
				conds = append(conds, clause.Gt{
					Column: exp.Column,
					Value:  exp.Value,
				})
			case "Lt":
				conds = append(conds, clause.Lt{
					Column: exp.Column,
					Value:  exp.Value,
				})
			case "Gte":
				conds = append(conds, clause.Gte{
					Column: exp.Column,
					Value:  exp.Value,
				})
			case "Lte":
				conds = append(conds, clause.Lte{
					Column: exp.Column,
					Value:  exp.Value,
				})
			}
		}
	}

	return conds
}

func (h gormPostgresqlDB) List(ctx context.Context, input *ListParams) *gorm.DB {
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

func initGormPostgresqlDB(options *GormConnectionOptions) (*gorm.DB, error) {
	timezone := "UTC"
	if options.Timezone != "" {
		timezone = options.Timezone
	}
	dsn := fmt.Sprintf(
		// "host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Asia/Ho_Chi_Minh",
		// "host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Local",
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable search_path=%s TimeZone=%s",
		options.Host,
		options.Port,
		options.Username,
		options.Password,
		options.Database,
		options.Schema,
		timezone,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		// DSN: "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), options.GormConfig)
	if err != nil {
		panic(err)
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
	// sqlDB.SetMaxIdleConns(10) // currently 2, maybe change in future
	// sqlDB.SetMaxOpenConns(100) // 0 unlimited
	// sqlDB.SetConnMaxIdleTime(time.Hour * 1)
	// sqlDB.SetConnMaxLifetime(time.Hour * 1)
	if options.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(options.MaxOpenConns)
	}

	zap.S().Infof("Successfully connected!")
	return db, nil
}

func pretreatmentConfig(enableLog bool) {
	config := new(gorm.Config)

	namingStrategy := schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	if enableLog {
		config.Logger = newLogger
	}

	config.NamingStrategy = namingStrategy
	config.PrepareStmt = true
}

type CRUDHelper interface {
	// Create creates a new record on database.
	// `input` must be a non-nil pointer of the model. e.g: `input := &model.User{}`
	Create(db *gorm.DB, input interface{}) error
	// View returns single record matching the given conditions.
	// `output` must be a non-nil pointer of the model. e.g: `output := new(model.User)`
	// Note: RecordNotFound error is returned when there is no record that matches the conditions
	View(db *gorm.DB, output interface{}, cond ...interface{}) error
	// List returns list of records retrievable after filter & pagination if given.
	// `output` must be a non-nil pointer of slice of the model. e.g: `data := []*model.User{}; db.List(dbconn, &data, nil, nil)`
	// `lq` can be nil, then no filter & pagination are applied
	// `count` can also be nil, then no extra query is executed to get the total count
	List(db *gorm.DB, output interface{}, lq *ListQueryCondition, count *int64) error
	// Update updates data of the records matching the given conditions.
	// `updates` could be a model struct or map[string]interface{}
	// Note: DB.Model must be provided in order to get the correct model/table
	Update(db *gorm.DB, updates interface{}, cond ...interface{}) error
	// Delete deletes record matching given conditions.
	// `cond` can be an instance of the model, then primary key will be used as the condition
	Delete(db *gorm.DB, cond ...interface{}) error
	// Exist checks whether there is record matching the given conditions.
	Exist(db *gorm.DB, cond ...interface{}) (bool, error)
	// CreateInBatches creates batch of new record on database.
	// `input` must be an array non-nil pointer of the model. e.g: `input := []*model.User`
	CreateInBatches(db *gorm.DB, input interface{}, batchSize int) error
	// DeletePermanently deletes record matching given conditions permanently.
	// `cond` can be an instance of the model, then primary key will be used as the condition
	DeletePermanently(db *gorm.DB, cond ...interface{}) error
}

type crudHelper struct {
	// Model must be set to a specific model instance. e.g: model.User{}
	Model interface{}
	// GDB holds previous DB instance that just executed the query
	GDB *gorm.DB
}

func NewCRUDHelper(gDB *gorm.DB, model interface{}) CRUDHelper {
	return &crudHelper{
		GDB:   gDB,
		Model: model,
	}
}

// Create creates a new record on database.
func (cdb *crudHelper) Create(db *gorm.DB, input interface{}) error {
	cdb.GDB = db.Create(input)
	return cdb.GDB.Error
}

// View returns single record matching the given conditions.
func (cdb *crudHelper) View(db *gorm.DB, output interface{}, cond ...interface{}) error {
	where := ParseCond(cond...)
	cdb.GDB = db.First(output, where...)
	return cdb.GDB.Error
}

// List returns list of records retrievable after filter & pagination if given.
func (cdb *crudHelper) List(db *gorm.DB, output interface{}, lq *ListQueryCondition, count *int64) error {
	if lq != nil {
		if lq.Filter != nil {
			db = db.Where(lq.Filter.SQL(), lq.Filter.Vars()...)
		}

		if lq.PerPage > 0 {
			db = db.Limit(lq.PerPage)
			if lq.Page > 1 {
				db = db.Offset(lq.Page*lq.PerPage - lq.PerPage)
			}
		}

		if lq.Sort != nil && len(lq.Sort) > 0 {
			// Note: It's up to whom using this package to validate the sort fields!
			db = db.Order(strings.Join(lq.Sort, ", "))
		}
	}

	cdb.GDB = db.Find(output)
	if err := cdb.GDB.Error; err != nil {
		return err
	}

	// Only count total records if requested
	if count != nil {
		if err := cdb.GDB.Limit(-1).Offset(-1).Count(count).Error; err != nil {
			return err
		}
	}

	return nil
}

// Update updates data of the records matching the given conditions.
func (cdb *crudHelper) Update(db *gorm.DB, updates interface{}, cond ...interface{}) error {
	db = db.Model(cdb.Model)
	if len(cond) > 0 {
		where := ParseCond(cond...)
		db = db.Where(where[0], where[1:]...)
	}
	cdb.GDB = db.Omit("id").Updates(updates)
	return cdb.GDB.Error
}

// Delete deletes record matching given conditions.
func (cdb *crudHelper) Delete(db *gorm.DB, cond ...interface{}) error {
	if len(cond) == 1 {
		val := reflect.ValueOf(cond[0])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			return db.Delete(cond[0]).Error
		}
	}
	where := ParseCond(cond...)
	cdb.GDB = db.Delete(cdb.Model, where...)
	return cdb.GDB.Error
}

// DeletePermanently deletes record matching given conditions permanently.
func (cdb *crudHelper) DeletePermanently(db *gorm.DB, cond ...interface{}) error {
	if len(cond) == 1 {
		val := reflect.ValueOf(cond[0])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			return db.Delete(cond[0]).Error
		}
	}
	where := ParseCond(cond...)
	cdb.GDB = db.Unscoped().Delete(cdb.Model, where...)
	return cdb.GDB.Error
}

// Exist checks whether there is record matching the given conditions.
func (cdb *crudHelper) Exist(db *gorm.DB, cond ...interface{}) (bool, error) {
	var count int64 = 0
	where := ParseCond(cond...)
	cdb.GDB = db.Model(cdb.Model).Where(where[0], where[1:]...).Count(&count)
	return count > 0, cdb.GDB.Error
}

// CreateInBatches creates batch of new record on database.
func (cdb *crudHelper) CreateInBatches(db *gorm.DB, input interface{}, batchSize int) error {
	cdb.GDB = db.CreateInBatches(input, batchSize)
	return cdb.GDB.Error
}

func ParseCondWithConfig(cfg gowhere.Config, cond ...interface{}) []interface{} {
	if len(cond) == 1 {
		switch c := cond[0].(type) {
		case map[string]interface{}, []interface{}:
			cond[0] = gowhere.WithConfig(cfg).Where(c)
		}

		if plan, ok := cond[0].(*gowhere.Plan); ok {
			return append([]interface{}{plan.SQL()}, plan.Vars()...)
		}
	}
	return cond
}

// ParseCond returns standard [sqlString, vars] format for query, powered by go_where package (with default config)
func ParseCond(cond ...interface{}) []interface{} {
	return ParseCondWithConfig(gowhere.DefaultConfig, cond...)
}

// ListQueryCondition holds data used for db queries
type ListQueryCondition struct {
	Filter  *gowhere.Plan
	Sort    []string
	Page    int
	PerPage int
}
