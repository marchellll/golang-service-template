package services

type EchoService interface {
	Echo(str string) string
}


type echoService struct {
}

func NewEchoService() EchoService {
	return &echoService{}
}

func (*echoService) Echo(str string) string {
	return str + " from service"
}