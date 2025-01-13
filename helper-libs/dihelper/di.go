package dihelper

import (
	"sync"

	"github.com/sarulabs/di"
)

type DIBuilder func() []di.Def

var (
	buildOnce           sync.Once
	builder             *di.Builder
	container           di.Container
	ConfigsBuilder      DIBuilder
	HelpersBuilder      DIBuilder
	RepositoriesBuilder DIBuilder
	AdaptersBuilder     DIBuilder
	UsecasesBuilder     DIBuilder
	FacadiesBuilder     DIBuilder
	APIsBuilder         DIBuilder
)

func BuildLibDIContainer() {
	buildOnce.Do(func() {
		builder, _ = di.NewBuilder()
		doBuild()
		container = builder.Build()
	})
}

func doBuild() {
	if err := buildConfigs(); err != nil {
		panic(err)
	}
	if err := buildHelpers(); err != nil {
		panic(err)
	}
	if err := buildRepositories(); err != nil {
		panic(err)
	}
	if err := buildAdapters(); err != nil {
		panic(err)
	}
	if err := buildUsecases(); err != nil {
		panic(err)
	}
	if err := buildFacadies(); err != nil {
		panic(err)
	}
	if err := buildAPIs(); err != nil {
		panic(err)
	}
}

func GetLibDependency(dependencyName string) interface{} {
	return container.Get(dependencyName)
}

func CleanDependency() error {
	return container.Clean()
}

func buildConfigs() error {
	// defs := []di.Def{}
	if ConfigsBuilder == nil {
		ConfigsBuilder = defaultBuilder
	}
	defs := ConfigsBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildHelpers() error {
	// defs := []di.Def{}
	if HelpersBuilder == nil {
		HelpersBuilder = defaultBuilder
	}
	defs := HelpersBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildRepositories() error {
	// defs := []di.Def{}
	if RepositoriesBuilder == nil {
		RepositoriesBuilder = defaultBuilder
	}
	defs := RepositoriesBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildAdapters() error {
	// defs := []di.Def{}
	if AdaptersBuilder == nil {
		AdaptersBuilder = defaultBuilder
	}
	defs := AdaptersBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildUsecases() error {
	// defs := []di.Def{}
	if UsecasesBuilder == nil {
		UsecasesBuilder = defaultBuilder
	}
	defs := UsecasesBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildFacadies() error {
	// defs := []di.Def{}
	if FacadiesBuilder == nil {
		FacadiesBuilder = defaultBuilder
	}
	defs := FacadiesBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildAPIs() error {
	// defs := []di.Def{}
	if APIsBuilder == nil {
		APIsBuilder = defaultBuilder
	}
	defs := APIsBuilder()
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func defaultBuilder() []di.Def {
	return []di.Def{}
}
