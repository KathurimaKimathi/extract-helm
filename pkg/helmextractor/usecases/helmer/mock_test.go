package helmer_test

import (
	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/infrastructure"
	fakeHelm "github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/infrastructure/services/helm/mock"
	helmUsecase "github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/usecases/helmer"
)

type Mocks struct {
	HelmMock *fakeHelm.HelmMock
}

func setupMocks() (helmUsecase.Helm, Mocks) {
	fakeHelm := fakeHelm.NewHelmMock()

	infra := infrastructure.NewInfrastructureInteractor(fakeHelm)
	useCase := helmUsecase.NewHelm(infra)

	mocks := Mocks{
		HelmMock: fakeHelm,
	}

	return *useCase, mocks
}
