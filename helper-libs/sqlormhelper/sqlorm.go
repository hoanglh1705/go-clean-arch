package sqlormhelper

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"go-clean-arch/helper-libs/commonhelper"
	"go-clean-arch/helper-libs/timehelper"
	"go-clean-arch/helper-libs/uuidhelper"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	DEFAULT_QUERY_LIMIT = 12
)

type (
	SqlGormDatabase interface {
		Open() *gorm.DB
		Close() error
		Begin() *gorm.DB
		Commit(tx *gorm.DB) *gorm.DB
		Rollback(tx *gorm.DB) *gorm.DB
		GetConn() (*gorm.DB, error)
		List(ctx context.Context, input *ListParams) *gorm.DB
	}

	GormConnectionOptions struct {
		Host               string
		Port               int
		Username           string
		Password           string
		Database           string
		Schema             string
		Timezone           string
		GormConfig         *gorm.Config
		UseTls             bool
		TlsMode            string
		TlsRootCACertFile  string
		TlsKeyFile         string
		TlsCertFile        string
		InsecureSkipVerify bool
		MaxOpenConns       int
	}

	// IdBaseEntity
	IdBaseEntity struct {
		Id               string         `gorm:"primaryKey"`
		CreatedTs        int64          `gorm:"column:created_ts;autoCreateTime:milli"`
		CreatedUser      string         `gorm:"column:created_user"`
		LastModifiedTs   int64          `gorm:"column:last_modified_ts;autoUpdateTime:milli"`
		LastModifiedUser string         `gorm:"column:last_modified_user"`
		DeletedAt        gorm.DeletedAt `gorm:"index"`
	}

	// UuidBaseEntity
	BaseEntity struct {
		CreatedTime      time.Time `gorm:"column:created_time;type:timestamptz;autoCreateTime"`
		CreatedUser      string    `gorm:"column:created_user;type:text"`
		LastModifiedTime time.Time `gorm:"column:last_modified_time;type:timestamptz;autoUpdateTime"`
		LastModifiedUser string    `gorm:"column:last_modified_user;type:text"`
	}

	BaseEntityWithId struct {
		Id               string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
		CreatedTime      time.Time      `gorm:"column:created_time;type:timestamptz;autoCreateTime"`
		CreatedUser      string         `gorm:"column:created_user;type:text"`
		LastModifiedTime time.Time      `gorm:"column:last_modified_time;type:timestamptz;autoUpdateTime"`
		LastModifiedUser string         `gorm:"column:last_modified_user;type:text"`
		DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at"`
	}

	BaseEntityTs struct {
		CreatedTs        int64          `gorm:"column:created_ts;autoCreateTime:milli"`
		CreatedUser      string         `gorm:"column:created_user"`
		LastModifiedTs   int64          `gorm:"column:last_modified_ts;autoUpdateTime:milli"`
		LastModifiedUser string         `gorm:"column:last_modified_user"`
		DeletedAt        gorm.DeletedAt `gorm:"index"`
	}

	BaseEntityTsWithId struct {
		Id               string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
		CreatedTs        int64          `gorm:"column:created_ts;autoCreateTime:milli"`
		CreatedUser      string         `gorm:"column:created_user"`
		LastModifiedTs   int64          `gorm:"column:last_modified_ts;autoUpdateTime:milli"`
		LastModifiedUser string         `gorm:"column:last_modified_user"`
		DeletedAt        gorm.DeletedAt `gorm:"index"`
	}

	ListParams struct {
		Model         interface{}
		Target        interface{}
		Limit         int
		Descending    bool
		OrderByColumn string
		LastQueryId   string
		WhereClauses  [][]string
	}

	GetParams struct {
		Model         interface{}
		Target        interface{}
		Id            string
		Descending    bool
		OrderByColumn string
		LastQueryId   string
		WhereClauses  [][]string
	}

	BaseListParams struct {
		Limit int
	}

	BaseUpdateInput struct {
		LastModifiedTime time.Time `gorm:"column:updated_date;autoUpdateTime"`
	}

	ConditionExpression struct {
		Column   string
		Value    string
		Operator string
	}
)

func (u *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context

	u.CreatedTime = timehelper.NewUTCTime()
	u.LastModifiedTime = timehelper.NewUTCTime()

	userId := ctx.Value(commonhelper.ContextKeyType_AppSubject)
	if userId != nil {
		userIdString := userId.(string)
		u.CreatedUser = userIdString
		u.LastModifiedUser = userIdString
	}

	return
}

func (u *BaseEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	u.LastModifiedTime = timehelper.NewUTCTime()

	userId := ctx.Value(commonhelper.ContextKeyType_AppSubject)
	if userId != nil {
		userIdString := userId.(string)
		u.LastModifiedUser = userIdString
	}

	return
}

func (u *BaseEntityWithId) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	u.Id = uuidhelper.NewUuidV7String()
	u.CreatedTime = timehelper.NewUTCTime()
	u.LastModifiedTime = timehelper.NewUTCTime()

	userId := ctx.Value(commonhelper.ContextKeyType_AppSubject)
	if userId != nil {
		userIdString := userId.(string)
		u.CreatedUser = userIdString
		u.LastModifiedUser = userIdString
	}

	return
}

func (u *BaseEntityWithId) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	u.LastModifiedTime = timehelper.NewUTCTime()

	userId := ctx.Value(commonhelper.ContextKeyType_AppSubject)
	if userId != nil {
		userIdString := userId.(string)
		u.LastModifiedUser = userIdString
	}

	return
}

type BaseColumnNameType string

const (
	BaseColumnNameType_Id               = "id"
	BaseColumnNameType_CreatedTime      = "created_time"
	BaseColumnNameType_CreatedUser      = "created_user"
	BaseColumnNameType_LastModifiedTime = "last_modified_time"
	BaseColumnNameType_LastModifiedUser = "last_modified_user"
	BaseColumnNameType_DeletedAt        = "deleted_at"
)

const (
	ExpressionOperator_Eq  = "Eq"
	ExpressionOperator_Gt  = "Gt"
	ExpressionOperator_Lt  = "Lt"
	ExpressionOperator_Gte = "Gte"
	ExpressionOperator_Lte = "Lte"
)

type Jsonb map[string]interface{}

func (j Jsonb) Value() (driver.Value, error) {
	if j == nil {
		j = map[string]interface{}{}
	}
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *Jsonb) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

func (Jsonb) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// use field.Tag, field.TagSettings gets field's tags
	// checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
