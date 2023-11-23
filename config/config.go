package config

type Config struct {
	AppName              string
	DBAddr               string
	DBUser               string
	DBPassword           string
	DBName               string
	OTLPExporterEndpoint string
	OTLPMetricExport     bool
}

func Load() Config {
	return Config{
		AppName:              "trace-app-golang",
		DBAddr:               "localhost:5432",
		DBUser:               "postgres",
		DBPassword:           "postgres",
		DBName:               "example",
		OTLPExporterEndpoint: "grafana-agent-flow.tempo.svc.cluster.local:4317",
		OTLPMetricExport:     false,
	}
}
