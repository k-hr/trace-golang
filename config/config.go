package config

type Config struct {
	AppName              string
	DBAddr               string
	DBUser               string
	DBPassword           string
	DBName               string
	OTLPExporterEndpoint string
}

func Load() Config {
	return Config{
		AppName:              "trace-app-golang",
		DBAddr:               "localhost:5432",
		DBUser:               "postgres",
		DBPassword:           "postgres",
		DBName:               "example",
		OTLPExporterEndpoint: "grafana-agent-traces.tempo.svc.cluster.local:4317",
	}
}
