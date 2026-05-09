package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Charset  string
}

type JWTConfig struct {
	Secret string
	Expire int
}

var AppConfig *Config

func InitConfig() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := filepath.Join(basepath, "..", "..", "configs")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found: %v, using defaults", err)
	}

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "root123")
	viper.SetDefault("database.dbname", "admin_system")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("jwt.secret", "admin-system-jwt-secret-key")
	viper.SetDefault("jwt.expire", 24)

	AppConfig = &Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Mode: viper.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			Charset:  viper.GetString("database.charset"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("jwt.secret"),
			Expire: viper.GetInt("jwt.expire"),
		},
	}

	log.Println("Config initialized successfully")
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}
