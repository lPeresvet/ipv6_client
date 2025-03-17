package service

type IfaceService struct{}

func NewIfaceService() *IfaceService {
	return &IfaceService{}
}

func (i IfaceService) GetIpv6Address(interfaceName string) (string, error) {
	return "some ipv6 address", nil
}
