package dihelper

import "github.com/sarulabs/di"

const (
	ConfigDIName string = "Config"
)

func BuildDIContainer() {
	initBuilder()
	BuildLibDIContainer()
}

func GetDependency(name string) interface{} {
	return GetLibDependency(name)
}

func initBuilder() {
	ConfigsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr, di.Def{
			Name:  ConfigDIName,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return nil, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		})

		return arr
	}
}
