package healthcheck

func (s *service) HealthcheckServices() (string, error) {
	return "service healty", nil
}
