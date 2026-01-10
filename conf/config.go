package conf

import (
	"fmt"
	"log"
	"mianshiba/pkg/lang/ternary"
	"mianshiba/pkg/logs"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config 应用程序配置结构
type Config struct {
	Host     string         `yaml:"host"`
	Port     int            `yaml:"port"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Hertz    HertzConfig    `yaml:"hertz"`
	Security SecurityConfig `yaml:"security"`
	API      APIConfig      `yaml:"api"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

type APIConfig struct {
	APIKeySecret string `yaml:"api_key_secret"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	DSN             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	DialTimeout  string `yaml:"dial_timeout"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// HertzConfig Hertz框架配置
type HertzConfig struct {
	LogLevel     string `yaml:"log_level"`
	LogPath      string `yaml:"log_path"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

// SecurityConfig 安全性配置
type SecurityConfig struct {
	JWTSecret        string `yaml:"jwt_secret"`
	JWTExpiration    string `yaml:"jwt_expiration"`
	JWTExpirationInt int64
	CORS             CORSConfig `yaml:"cors"`
}

func (c *Config) ExpandEnv() {
	c.Redis.Password = expandEnvVar(c.Redis.Password)
	c.Security.JWTSecret = expandEnvVar(c.Security.JWTSecret)
	c.API.APIKeySecret = expandEnvVar(c.API.APIKeySecret)
}

// Global 全局配置实例
var Global Config

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &Global)
	if err != nil {
		return err
	}

	err = loadEnv()
	if err != nil {
		return err
	}

	Global.ExpandEnv()

	expiration, err := time.ParseDuration(Global.Security.JWTExpiration)
	if err != nil {
		return err
	}
	Global.Security.JWTExpirationInt = int64(expiration.Seconds())

	log.Println("配置加载成功")
	return nil
}

func loadEnv() (err error) {
	appEnv := os.Getenv("APP_ENV")
	fileName := ternary.IFElse(appEnv == "", ".env", ".env."+appEnv)

	logs.Infof("load env file: %s", fileName)

	err = godotenv.Load(fileName)
	if err != nil {
		return fmt.Errorf("load env file(%s) failed, err=%w", fileName, err)
	}

	return err
}

// expandEnvVar 展开字符串中的环境变量引用
// 支持 ${VAR_NAME} 和 $VAR_NAME 两种语法
func expandEnvVar(s string) string {
	if s == "" {
		return s
	}

	// 匹配 ${VAR_NAME} 或 $VAR_NAME
	re := regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

	result := re.ReplaceAllStringFunc(s, func(match string) string {
		// 提取变量名
		varName := ""
		if strings.HasPrefix(match, "${") {
			// ${VAR_NAME} 格式
			varName = match[2 : len(match)-1]
		} else {
			// $VAR_NAME 格式
			varName = match[1:]
		}

		// 获取环境变量值
		value := os.Getenv(varName)
		if value != "" {
			return value
		}

		// 如果环境变量不存在，保持原样
		return match
	})

	return result
}
