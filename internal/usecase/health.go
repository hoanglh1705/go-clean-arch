package usecase

type (
	HealthUsecase interface {
		Liveness() error
		Readiness() error
	}

	healthUsecase struct {
	}
)

func NewHealthUsecase() HealthUsecase {
	return &healthUsecase{}
}

func (u *healthUsecase) Liveness() error {
	return nil
}

func (u *healthUsecase) Readiness() error {
	return nil
}
