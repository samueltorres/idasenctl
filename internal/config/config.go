package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	ErrDeskNotExists = errors.New("desk not exists")
)

type Config struct {
	Desks       map[string]Desk `yaml:"desks"`
	DefaultDesk string          `yaml:"defaultDesk"`
}

type Desk struct {
	Name    string            `yaml:"name"`
	Address string            `yaml:"address"`
	Presets map[string]Preset `yaml:"presets"`
}

type Preset struct {
	Name   string  `yaml:"name"`
	Height float32 `yaml:"height"`
}

type ConfigManager struct {
	configFile string
	config     *Config
}

func NewConfigManager(configFile string) (*ConfigManager, error) {
	config, err := readConfigFromFile(configFile)
	if err != nil {
		return nil, err
	}

	if config.Desks == nil {
		config.Desks = make(map[string]Desk)
	}

	return &ConfigManager{
		config:     &config,
		configFile: configFile,
	}, nil
}

func (cm *ConfigManager) GetDesk(name string) (Desk, error) {
	if d, ok := cm.config.Desks[name]; ok {
		return d, nil
	}

	return Desk{}, ErrDeskNotExists
}

func (cm *ConfigManager) SetDesk(desk Desk) error {
	cm.config.Desks[desk.Name] = desk
	return cm.storeConfig()
}

func (cm *ConfigManager) SetDefaultDesk(name string) error {
	cm.config.DefaultDesk = name
	return cm.storeConfig()
}

func (cm *ConfigManager) GetDefaultDesk() string {
	return cm.config.DefaultDesk
}

func (cm *ConfigManager) SetDeskPreset(deskName string, presetName string, height float32) error {
	d, ok := cm.config.Desks[deskName]
	if !ok {
		return ErrDeskNotExists
	}

	if d.Presets == nil {
		d.Presets = make(map[string]Preset)
	}
	d.Presets[strings.ToLower(presetName)] = Preset{
		Name:   presetName,
		Height: height,
	}

	return cm.storeConfig()
}

func (cm *ConfigManager) DeleteDeskPreset(deskName string, presetName string) error {
	d, ok := cm.config.Desks[deskName]
	if !ok {
		return ErrDeskNotExists
	}

	if d.Presets == nil {
		d.Presets = make(map[string]Preset)
	}

	delete(d.Presets, strings.ToLower(presetName))

	return cm.storeConfig()
}

func (cm *ConfigManager) storeConfig() error {
	f, err := os.OpenFile(cm.configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "could not open config file")
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cm.config)
	if err != nil {
		return errors.Wrap(err, "could not save config file")
	}

	return nil
}

func readConfigFromFile(configFile string) (Config, error) {
	f, err := os.OpenFile(configFile, os.O_CREATE, 0644)
	if err != nil {
		return Config{}, errors.Wrap(err, "could not open config file")
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, errors.Wrap(err, "could not read config file")
	}

	if len(b) == 0 {
		return Config{}, nil
	}

	var cfg Config
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "could not parse config file")
	}

	return cfg, nil
}
