package config

import (
	"bytes"
	"fmt"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/viper"
	"io/ioutil"
	"os"
)

func ReadConfig(configPath string) (registry configcore.Registry, serverConfig Application) {
	registry = viper.New()
	registry.SetConfigType("yaml")

	f, err := os.Open(configPath)
	if err != nil {
		panic("Cannot read config file: " + err.Error())
	}
	defer f.Close()

	if err = registry.ReadConfig(f); err != nil {
		panic("Cannot read config file: " + err.Error())
	}
	// server config
	if err := registry.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	return
}

func LoadSecureConfigFiles(registry configcore.Registry,
	serverConfig Application) []SecureConfigDetail {

	var result []SecureConfigDetail
	loadedSecuredConfigs := make(map[string]SecureConfigDetail)

	for _, coin := range serverConfig.AvailableCoins {
		coinConfig := registry.ValueOf(coin)
		var cc struct {
			SecureConfigPath string
		}
		if err := coinConfig.Unmarshal(&cc); err != nil {
			panic(err)
		}
		if cc.SecureConfigPath == "" {
			panic(fmt.Sprintf("%s: invalid secure config path value", coin))
		}
		if cfg, ok := loadedSecuredConfigs[cc.SecureConfigPath]; ok {
			detail := SecureConfigDetail{
				Coin:     coin,
				FilePath: cc.SecureConfigPath,
				Data:     cfg.Data,
			}
			result = append(result, detail)
			continue
		}
		f, err := os.Open(cc.SecureConfigPath)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
		if len(data) < 10 {
			panic("empty file")
		}
		detail := SecureConfigDetail{
			Coin:     coin,
			FilePath: cc.SecureConfigPath,
			Data:     data,
		}
		result = append(result, detail)
		loadedSecuredConfigs[cc.SecureConfigPath] = detail
	}
	return result
}

func IsAnyConfigSecured(configs []SecureConfigDetail) bool {
	for _, cfg := range configs {
		if bytes.Compare(cfg.Data[:10], lockedKey) == 0 {
			return true
		}
	}
	return false
}

func LoadSecureConfigs(configs []SecureConfigDetail,
	unlockHandler UnlockHandler) (secureConfigs map[string]SecureConfig, err error) {

	secureConfigs = make(map[string]SecureConfig)

	for _, cfg := range configs {
		if bytes.Compare(cfg.Data[:10], lockedKey) == 0 {
			// unlock
			data := cfg.Data[10:]
			for {
				unlockedData, err := unlockHandler(data)
				if err != nil {
					fmt.Println(fmt.Sprintf("error on unlock. err: %v", err))
					continue
				}
				cfg.Data = unlockedData
				break
			}
		}
		reader := bytes.NewReader(cfg.Data)
		secureConf := viper.New()
		secureConf.SetConfigType("yaml")
		if err := secureConf.ReadConfig(reader); err != nil {
			return nil, err
		}

		var sc SecureConfig
		if err := secureConf.Unmarshal(&sc); err != nil {
			return nil, err
		}
		sc.detail = cfg
		secureConfigs[cfg.Coin] = sc
	}
	return secureConfigs, err
}
