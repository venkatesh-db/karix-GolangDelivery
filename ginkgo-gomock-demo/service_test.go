package ginkgo_gomock_demo_test

import (
	"errors"
	"testing"

	demo "ginkgo-gomock-demo"
	"ginkgo-gomock-demo/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

install two packages
go get github.com/onsi/ginkgo/v2
go get github.com/onsi/gomega



*/

func TestService(t *testing.T) {
	
	RegisterFailHandler(Fail)
//	RunSpecs(t, "UserService Suite")
	RunSpecs(t, "Jamesbond Suite")
}

var _ = Describe("Jamesbond", func() {

	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockUserRepo
		svc      *demo.UserService
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockUserRepo(ctrl)
		svc = &demo.UserService{Repo: mockRepo}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("RobPikeSmile()", func() {

		It("returns Rob Pike smile message", func() {

			mockRepo.EXPECT().
				RobPike(":-)").
				Return(":-D", nil)
			msg, err := svc.RobPikeSmile(":-)")


			Expect(err).To(BeNil())
			Expect(msg).To(Equal("Rob Pike says: :-D"))
		})

	})

})

var _ = Describe("UserService", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockUserRepo
		svc      *demo.UserService
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockUserRepo(ctrl)
		svc = &demo.UserService{Repo: mockRepo}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Welcome()", func() {

		It("returns welcome message", func() {
			mockRepo.EXPECT().
				GetName(1).
				Return("Venkatesh1", nil)

			msg, err := svc.Welcome(1)

			Expect(err).To(BeNil())
			Expect(msg).To(Equal("Welcome Venkatesh"))
		})

		It("returns error", func() {
			mockRepo.EXPECT().
				GetName(2).
				Return("", errors.New("not found"))

			msg, err := svc.Welcome(2)

			Expect(err).To(HaveOccurred())
			Expect(msg).To(Equal(""))
		})
	})
})
