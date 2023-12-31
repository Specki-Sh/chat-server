package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"chat-server/internal/service"
	"chat-server/pkg/db"
	"chat-server/pkg/redis"
	"chat-server/pkg/server"
)

type Config struct {
	Server struct {
		Host           string        `mapstructure:"host"`
		Port           int           `mapstructure:"port"`
		MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
		ReadTimeout    time.Duration `mapstructure:"read_timeout"`
		WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	} `mapstructure:"server"`
	DB struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	Token struct {
		AccessKeys struct {
			PublicKeyPath  string `mapstructure:"public_key_path"`
			PublicKey      *rsa.PublicKey
			PrivateKeyPath string `mapstructure:"private_key_path"`
			PrivateKey     *rsa.PrivateKey
		} `mapstructure:"access_keys"`
		RefreshKeys struct {
			PublicKeyPath  string `mapstructure:"public_key_path"`
			PublicKey      *rsa.PublicKey
			PrivateKeyPath string `mapstructure:"private_key_path"`
			PrivateKey     *rsa.PrivateKey
		} `mapstructure:"refresh_keys"`
		AccessExpiration  time.Duration `mapstructure:"access_expiration"`
		RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
	} `mapstructure:"token"`
}

func (c *Config) Parse() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Config.Parse: %w", err)
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("Config.Parse: %w", err)
	}

	c.Token.AccessKeys.PublicKey, err = readPublicKeyFile(c.Token.AccessKeys.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("Config.Parse: %w", err)
	}
	c.Token.AccessKeys.PrivateKey, err = readPrivateKeyFile(c.Token.AccessKeys.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("Config.Parse: %w", err)
	}

	c.Token.RefreshKeys.PublicKey, err = readPublicKeyFile(c.Token.RefreshKeys.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("Config.Parse: %w", err)
	}
	c.Token.RefreshKeys.PrivateKey, err = readPrivateKeyFile(c.Token.RefreshKeys.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("Config.Parse: %w", err)
	}

	return nil
}

func (c *Config) GetServerConfig() *server.Config {
	return &server.Config{
		Addr:           fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port),
		MaxHeaderBytes: c.Server.MaxHeaderBytes,
		ReadTimeout:    c.Server.ReadTimeout,
		WriteTimeout:   c.Server.WriteTimeout,
	}
}

func (c *Config) GetDBConfig() *db.Config {
	return &db.Config{
		Host:     c.DB.Host,
		Port:     strconv.Itoa(c.DB.Port),
		Username: c.DB.Username,
		Password: c.DB.Password,
		DBName:   c.DB.DBName,
		SSLMode:  c.DB.SSLMode,
	}
}

func (c *Config) GetRedisConfig() *redis.Config {
	return &redis.Config{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	}
}

func (c *Config) GetTSConfig() *service.TSConfig {
	return &service.TSConfig{
		AccessKeys: &service.KeyPair{
			PrivKey: c.Token.AccessKeys.PrivateKey,
			PubKey:  c.Token.AccessKeys.PublicKey,
		},
		RefreshKeys: &service.KeyPair{
			PrivKey: c.Token.RefreshKeys.PrivateKey,
			PubKey:  c.Token.RefreshKeys.PublicKey,
		},
		AccessExpiration:  &c.Token.AccessExpiration,
		RefreshExpiration: &c.Token.RefreshExpiration,
	}
}

func readPublicKeyFile(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readPublicKeyFile: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("readPublicKeyFile: no PEM data found")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("readPublicKeyFile: %w", err)
	}
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("readPublicKeyFile: not an RSA key")
	}
	return pub, nil
}

func readPrivateKeyFile(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readPrivateKeyFile: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("readPrivateKeyFile: no PEM data found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("readPrivateKeyFile: %w", err)
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("readPrivateKeyFile: not an RSA key")
	}
	return priv, nil
}
