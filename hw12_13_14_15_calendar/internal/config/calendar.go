package config

type CalendarConfig struct {
	HTTP    HTTPConfig
	RPC     RPCConfig
	Storage StorageConfig
	Log     LogConfig
}

func NewCalendarConfig() *CalendarConfig {
	return &CalendarConfig{}
}
