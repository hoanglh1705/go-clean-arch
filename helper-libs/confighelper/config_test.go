package confighelper

import (
	"go-clean-arch/helper-libs/loghelper"
	"testing"
)

var defaultConfig = []byte(`
environment: dev
app: pm-helper-libs
`)

type (
	Config struct {
		Environment string `mapstructure:"environment"`
		App         string `mapstructure:"app"`
	}
)

func LoadConfig() (*Config, error) {
	var cfg = &Config{}
	err := Load(cfg, defaultConfig)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func TestLoadConfigSuccess(t *testing.T) {
	_ = loghelper.InitZap("testing", "dev", map[string]string{})
	cfg, _ := LoadConfig()
	got := cfg.Environment
	want := "dev"
	loghelper.Logger.Infof("Test Load config success")
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestLoadConfigFailure(t *testing.T) {
	cfg, _ := LoadConfig()
	got := cfg.Environment
	want := "devv"
	if got == want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
