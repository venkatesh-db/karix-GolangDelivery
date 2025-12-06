package ginkgo_gomock_demo

/*

mockgen -source=service.go -destination=mocks/mock_userrepo.go -package=mocks


*/

type UserRepo interface {
	GetName(id int) (string, error)
	RobPike(smile string) (string, error)
}

type UserService struct {
	Repo UserRepo
}

func (s *UserService) Welcome(id int) (string, error) {
	name, err := s.Repo.GetName(id)
	if err != nil {
		return "", err
	}
	return "Welcome " + name, nil
}

func (s *UserService) RobPikeSmile(smile string) (string, error) {

	result, err := s.Repo.RobPike(smile)
	if err != nil {
		return "", err
	}
	return "Rob Pike says: " + result, nil
}
