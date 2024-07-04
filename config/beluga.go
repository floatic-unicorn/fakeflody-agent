package config

var (
	devToken         = "eyJhbGciOiJIUzUxMiJ9.eyJ0eXBlIjoiYWNjZXNzVG9rZW4iLCJhY2NvdW50SWQiOiJyYktEdkU4TyIsImVtYWlsIjoic3ZjQGZsb2F0aWMuaW8iLCJyb2xlIjoiSU5URVJOQUxfU1ZDIiwiZXhwIjozMjYyMDYyOTYwfQ.sqvC0i6EkKAP7xfijvbnKXwvlrzD-nIFGsOs0Oo5gHj_UAnNaOSCBGAzvg-1ypZkTzRbTu8_CVyLriMgtQQRHw"
	robotInfoService = EndpointConfig{
		Url: "https://dev.beluga.floatic.io/robotInfo",
	}
)

func NewBelugaConfig() BelugaConfig {
	return BelugaConfig{
		RobotInfoService: robotInfoService,
		Token:            devToken,
	}
}

type BelugaConfig struct {
	RobotInfoService EndpointConfig
	Token            string
}

type EndpointConfig struct {
	Url string
}
