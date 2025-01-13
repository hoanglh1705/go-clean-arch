package copyhelper

import (
	"go-clean-arch/helper-libs/timehelper"

	"github.com/jinzhu/copier"
)

type (
	CopyOptions struct {
		Timezone timehelper.Timezone
	}
)

type (
	PbConverter interface {
		FromPb(to interface{}, from interface{})
		ToPb(to interface{}, from interface{})
	}

	ModelConverter interface {
		FromModel(to interface{}, from interface{})
		ToModel(to interface{}, from interface{})
	}

	EntityConverter interface {
		FromEntity(to interface{}, from interface{})
		ToEntity(to interface{}, from interface{})
	}

	ICoreEntityConverter interface {
		FromModel(to interface{}, from interface{})
		ToModel(to interface{}, from interface{})
	}

	AdapterConverter interface {
		FromAdapter(to interface{}, from interface{})
		ToAdapter(to interface{}, from interface{})
	}

	ObjectCopier interface {
		Copy(to interface{}, from interface{})
	}

	IgnoreEmptyObjectCopier interface {
		Copy(to interface{}, from interface{})
	}

	modelConverter struct {
	}

	entityConverter struct {
	}

	adapterConverter struct {
	}

	objectCopier struct {
	}

	ignoreEmptyObjectCopier struct {
	}
)

func NewModelConverter() ModelConverter {
	return &modelConverter{}
}

func (h *modelConverter) FromModel(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func (h *modelConverter) ToModel(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func NewEntityConverter() EntityConverter {
	return &entityConverter{}
}

func (h *entityConverter) FromEntity(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func (h *entityConverter) ToEntity(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func NewAdapterConverter() AdapterConverter {
	return &adapterConverter{}
}

func (h *adapterConverter) FromAdapter(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func (h *adapterConverter) ToAdapter(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func NewObjectCopier() ObjectCopier {
	return &objectCopier{}
}

func (h *objectCopier) Copy(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func NewIgnoreEmptyObjectCopier() ObjectCopier {
	return &ignoreEmptyObjectCopier{}
}

func (h *ignoreEmptyObjectCopier) Copy(to interface{}, from interface{}) {
	_ = copier.CopyWithOption(to, from, copier.Option{
		IgnoreEmpty: true,
	})
}
