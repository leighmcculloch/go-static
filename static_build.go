package static

func (s *Static) Build() error {
	for path := range s.pages {
		err := s.BuildPage(path)
		if err != nil {
			return err
		}
	}
	return nil
}
