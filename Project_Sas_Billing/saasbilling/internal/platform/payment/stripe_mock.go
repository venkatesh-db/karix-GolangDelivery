package payment

import "errors"

type GatewayResponse struct {
	Success      bool
	RefID        string
	ErrorMessage string
}

type Gatway interface {
	Charge(amount int64, currency, source, description string) (*GatewayResponse, error)
	Refund(chargeID string, amount int64) (*GatewayResponse, error)
}

type StripeMock struct {
	ShouldFail bool
}

func (s *StripeMock) Charge(amount int64, currency, source, description string) (*GatewayResponse, error) {
	if s.ShouldFail {
		return &GatewayResponse{
			Success:      false,
			ErrorMessage: "Charge failed",
		}, errors.New("charge failed")
	}

	return &GatewayResponse{
		Success: true,
		RefID:   "ch_mocked_12345",
	}, nil
}

func (s *StripeMock) Refund(chargeID string, amount int64) (*GatewayResponse, error) {

	if s.ShouldFail {
		return &GatewayResponse{
			Success:      false,
			ErrorMessage: "Refund failed",
		}, errors.New("refund failed")
	}

	return &GatewayResponse{
		Success: true,
		RefID:   "re_mocked_12345",
	}, nil
}
