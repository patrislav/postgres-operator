package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.bus.zalan.do/acid/postgres-operator/pkg/spec"
	"github.com/kelseyhightower/envconfig"
)

type TPR struct {
	ReadyWaitInterval time.Duration `split_words:"true" default:"3s"`
	ReadyWaitTimeout  time.Duration `split_words:"true" default:"30s"`
	ResyncPeriod      time.Duration `split_words:"true" default:"5m"`
}

type Resources struct {
	ResyncPeriodPod        time.Duration `split_words:"true" default:"5m"`
	ResourceCheckInterval  time.Duration `split_words:"true" default:"3s"`
	ResourceCheckTimeout   time.Duration `split_words:"true" default:"10m"`
	PodLabelWaitTimeout    time.Duration `split_words:"true" default:"10m"`
	PodDeletionWaitTimeout time.Duration `split_words:"true" default:"10m"`
}

type Auth struct {
	PamRoleName                   string              `split_words:"true" default:"zalandos"`
	PamConfiguration              string              `split_words:"true" default:"https://info.example.com/oauth2/tokeninfo?access_token= uid realm=/employees"`
	TeamsAPIUrl                   string              `envconfig:"teams_api_url" default:"https://teams.example.com/api/"`
	OAuthTokenSecretName          spec.NamespacedName `envconfig:"oauth_token_secret_name" default:"postgresql-operator"`
	InfrastructureRolesSecretName spec.NamespacedName `split_words:"true"`
	SuperUsername                 string              `split_words:"true" default:"postgres"`
	ReplicationUsername           string              `split_words:"true" default:"replication"`
}

type Config struct {
	TPR
	Resources
	Auth
	EtcdHost           string `split_words:"true" default:"etcd-client.default.svc.cluster.local:2379"`
	DockerImage        string `split_words:"true" default:"registry.opensource.zalan.do/acid/spilo-9.6:1.2-p12"`
	ServiceAccountName string `split_words:"true" default:"operator"`
	DbHostedZone       string `split_words:"true" default:"db.example.com"`
	EtcdScope          string `split_words:"true" default:"service"`
	WALES3Bucket       string `envconfig:"wal_s3_bucket"`
	KubeIAMRole        string `envconfig:"kube_iam_role"`
	DebugLogging       bool   `split_words:"true" default:"false"`
	DNSNameFormat      string `envconfig:"dns_name_format" default:"%s.%s.%s"`
}

func LoadFromEnv() *Config {
	//TODO: maybe we should use ConfigMaps( https://kubernetes.io/docs/tasks/configure-pod-container/configmap/ ) instead?

	var cfg Config
	err := envconfig.Process("PGOP", &cfg)
	if err != nil {
		panic(fmt.Errorf("Can't read config: %v", err))
	}
	cfg.EtcdScope = strings.Trim(cfg.EtcdScope, "/")

	return &cfg
}

func (c Config) MustMarshal() string {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}

	return string(b)
}