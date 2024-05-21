package services

type HealthService interface {
	Healthcheck(str string) string
}

type healthService struct {
}

func NewHealthService() HealthService {
	return &healthService{}
}

func (*healthService) Healthcheck(str string) string {
	return str + " from service"
}
