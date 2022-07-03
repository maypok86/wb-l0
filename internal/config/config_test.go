package config

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	type env struct {
		environment      string
		httpHost         string
		httpPort         string
		postgresHost     string
		postgresPort     string
		postgresDBName   string
		postgresUser     string
		postgresPassword string
		stanHost         string
		stanPort         string
		stanClusterID    string
		stanClientID     string
		stanChannel      string
	}

	setEnv := func(env env) {
		require.NoError(t, os.Setenv("ENVIRONMENT", env.environment))
		require.NoError(t, os.Setenv("HTTP_HOST", env.httpHost))
		require.NoError(t, os.Setenv("HTTP_PORT", env.httpPort))
		require.NoError(t, os.Setenv("POSTGRES_HOST", env.postgresHost))
		require.NoError(t, os.Setenv("POSTGRES_PORT", env.postgresPort))
		require.NoError(t, os.Setenv("POSTGRES_DBNAME", env.postgresDBName))
		require.NoError(t, os.Setenv("POSTGRES_USER", env.postgresUser))
		require.NoError(t, os.Setenv("POSTGRES_PASSWORD", env.postgresPassword))
		require.NoError(t, os.Setenv("STAN_HOST", env.stanHost))
		require.NoError(t, os.Setenv("STAN_PORT", env.stanPort))
		require.NoError(t, os.Setenv("STAN_CLUSTER_ID", env.stanClusterID))
		require.NoError(t, os.Setenv("STAN_CLIENT_ID", env.stanClientID))
		require.NoError(t, os.Setenv("STAN_CHANNEL", env.stanChannel))
	}

	tests := []struct {
		name string
		env  env
		want *Config
	}{
		{
			name: "test config",
			env: env{
				environment:      "test",
				httpHost:         "0.0.0.0",
				httpPort:         "8080",
				postgresHost:     "postgres",
				postgresPort:     "5431",
				postgresDBName:   "test_wb-l0",
				postgresUser:     "test_wb-l0",
				postgresPassword: "test",
				stanHost:         "0.0.0.0",
				stanPort:         "4222",
				stanClusterID:    "test-cluster",
				stanClientID:     "test-client",
				stanChannel:      "test",
			},
			want: &Config{
				Environment: test,
				HTTP: HTTP{
					Host:           "0.0.0.0",
					Port:           "8080",
					MaxHeaderBytes: 1,
					ReadTimeout:    10 * time.Second,
					WriteTimeout:   10 * time.Second,
				},
				Postgres: Postgres{
					Host:     "postgres",
					Port:     "5431",
					DBName:   "test_wb-l0",
					User:     "test_wb-l0",
					Password: "test",
					SSLMode:  "disable",
				},
				STAN: STAN{
					Host:      "0.0.0.0",
					Port:      "4222",
					ClusterID: "test-cluster",
					ClientID:  "test-client",
					Channel:   "test",
				},
				Logger: Logger{
					Level: "info",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.env)

			got := Get()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
