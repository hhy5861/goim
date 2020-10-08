package comet

type (
	HelloService struct {
	}
)

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (svc *HelloService) Action() {

}
