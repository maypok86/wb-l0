package nats

type Config struct {
	Host      string
	Port      string
	ClusterID string
	ClientID  string
}

func NewConfig(host string, port string, clusterID string, clientID string) Config {
	return Config{
		Host:      host,
		Port:      port,
		ClusterID: clusterID,
		ClientID:  clientID,
	}
}
