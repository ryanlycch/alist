package bootstrap

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alist-org/alist/v3/cmd/flags"
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

func InitConfig() {
	log.Infof("reading config file: %s", flags.Config)
	if !utils.Exists(flags.Config) {
		log.Infof("config file not exists, creating default config file")
		_, err := utils.CreateNestedFile(flags.Config)
		if err != nil {
			log.Fatalf("failed to create config file: %+v", err)
		}
		conf.Conf = conf.DefaultConfig()
		if !utils.WriteToJson(flags.Config, conf.Conf) {
			log.Fatalf("failed to create default config file")
		}
	} else {
		configBytes, err := ioutil.ReadFile(flags.Config)
		if err != nil {
			log.Fatalf("reading config file error: %+v", err)
		}
		conf.Conf = conf.DefaultConfig()
		err = utils.Json.Unmarshal(configBytes, conf.Conf)
		if err != nil {
			log.Fatalf("load config error: %+v", err)
		}
		// update config.json struct
		confBody, err := utils.Json.MarshalIndent(conf.Conf, "", "  ")
		if err != nil {
			log.Fatalf("marshal config error: %+v", err)
		}
		err = ioutil.WriteFile(flags.Config, confBody, 0777)
		if err != nil {
			log.Fatalf("update config struct error: %+v", err)
		}
	}
	if !conf.Conf.Force {
		confFromEnv()
	}
	// convert abs path
	var absPath string
	var err error
	if !filepath.IsAbs(conf.Conf.TempDir) {
		absPath, err = filepath.Abs(conf.Conf.TempDir)
		if err != nil {
			log.Fatalf("get abs path error: %+v", err)
		}
	}
	conf.Conf.TempDir = absPath
	err = os.RemoveAll(filepath.Join(conf.Conf.TempDir))
	if err != nil {
		log.Errorln("failed delete temp file:", err)
	}
	err = os.MkdirAll(conf.Conf.TempDir, 0700)
	if err != nil {
		log.Fatalf("create temp dir error: %+v", err)
	}
	log.Debugf("config: %+v", conf.Conf)
}

func confFromEnv() {
	prefix := "ALIST_"
	if flags.NoPrefix {
		prefix = ""
	}
	log.Infof("load config from env with prefix: %s", prefix)
	if err := env.Parse(conf.Conf, env.Options{
		Prefix: prefix,
	}); err != nil {
		log.Fatalf("load config from env error: %+v", err)
	}
}
