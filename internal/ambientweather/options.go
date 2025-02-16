package ambientweather

type Option func(s *Server)

func WithVersion(version string) Option {
	return func(s *Server) {
		s.version = version
	}
}

func WithUserAgent(userAgent string) Option {
	return func(s *Server) {
		s.userAgent = userAgent
	}
}
