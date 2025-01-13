package diregistry

import (
	"go-clean-arch/config"
	"go-clean-arch/helper-libs/copyhelper"
	"go-clean-arch/helper-libs/dihelper"
	"go-clean-arch/helper-libs/sqlormhelper"
	v1 "go-clean-arch/internal/api/v1"
	"go-clean-arch/internal/usecase"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

// DI Path
const (
	// Redis
	CacheHelperDIName       string = "RedisCacheHelper"
	RedisClientHelperDIName string = "RedisClientHelper"

	// Config
	ConfigDIName string = "Config"

	// Helper
	PbConverterDIName      string = "PbConverter"
	ModelConverterDIName   string = "ModelConverter"
	AdapterConverterDIName string = "AdapterConverter"
	SqlGormHelperDIName    string = "SqlGormHelper"

	DataBaseDIName string = "Database"

	// Repository
	BaseRepositoryDIName string = "BaseRepository"

	//Usecase
	HealthUsecaseDIName string = "HealthUsecase"

	// Api
	ApiServerV1DIName string = "ApiServerV1"
)

func BuildDIContainer() {
	initBuilder()
	dihelper.BuildLibDIContainer()
}

func GetDependency(name string) interface{} {
	return dihelper.GetLibDependency(name)
}

func initBuilder() {
	dihelper.ConfigsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr, di.Def{
			Name:  ConfigDIName,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg, err := config.Load()
				return cfg, err
			},
			Close: func(obj interface{}) error {
				return nil
			},
		})

		return arr
	}

	dihelper.HelpersBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr, di.Def{
			Name:  SqlGormHelperDIName,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDIName).(*config.Config)
				return sqlormhelper.NewGormPostgresqlDB(&sqlormhelper.GormConnectionOptions{
					Host:       cfg.Database.Host,
					Port:       int(cfg.Database.Port),
					Username:   cfg.Database.Username,
					Password:   cfg.Database.Password,
					Database:   cfg.Database.Database,
					Schema:     cfg.Database.SearchPath,
					GormConfig: &gorm.Config{},
				}), nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		}, di.Def{
			Name:  ModelConverterDIName,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return copyhelper.NewModelConverter(), nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		})
		return arr
	}

	dihelper.RepositoriesBuilder = func() []di.Def {
		arr := []di.Def{}
		// arr = append(arr,
		// 	di.Def{
		// 		Name:  BaseRepositoryDIName,
		// 		Scope: di.App,
		// 		Build: func(ctn di.Container) (interface{}, error) {
		// 			sql := ctn.Get(SqlGormHelperDIName).(sqlormhelper.SqlGormDatabase)
		// 			baseRepository := repository.NewBaseRepository(sql)
		// 			return baseRepository, nil
		// 		},
		// 		Close: func(obj interface{}) error {
		// 			return nil
		// 		},
		// 	},
		// )

		return arr
	}

	dihelper.UsecasesBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  HealthUsecaseDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (interface{}, error) {
					return usecase.NewHealthUsecase(), nil
				},
				Close: func(obj interface{}) error {
					return nil
				},
			},
		)

		return arr
	}

	dihelper.APIsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  ApiServerV1DIName,
				Scope: di.App,
				Build: func(ctn di.Container) (interface{}, error) {
					healthUsecase := ctn.Get(HealthUsecaseDIName).(usecase.HealthUsecase)
					return v1.NewAPIServer(healthUsecase), nil
				},
				Close: func(obj interface{}) error {
					return nil
				},
			},
		)
		return arr
	}
}
