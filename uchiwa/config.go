package uchiwa

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/palourde/logger"
)

// Config struct contains []SensuConfig and UchiwaConfig structs
type Config struct {
	Dashboard *GlobalConfig `json:",omitempty"`
	Sensu     []SensuConfig
	Uchiwa    GlobalConfig
}

// SensuConfig struct contains conf about a Sensu API
type SensuConfig struct {
	Name     string
	Host     string
	Port     int
	Ssl      bool
	Insecure bool
	URL      string
	User     string
	Path     string
	Pass     string
	Timeout  int
}

// GlobalConfig struct contains conf about Uchiwa
type GlobalConfig struct {
	Host    string
	Port    int
	Refresh int
	Pass    string
	User    string
	Db      Db
	Github  Github
	Auth    string
}

// Db struct contains the SQL driver configuration
type Db struct {
	Driver string
	Scheme string
}

// Github struct contains the GitHub driver configuration
type Github struct {
	ClientID     string
	ClientSecret string
	Roles        GithubRoles
	Server       string
}

// GithubRoles contains the roles of each GitHub team
type GithubRoles struct {
	Guests    []string
	Operators []string
}

func (c *Config) initSensu() {
	for i, api := range c.Sensu {
		prot := "http"
		if api.Name == "" {
			logger.Warningf("Sensu API %s has no name property. Generating random one...", api.URL)
			c.Sensu[i].Name = fmt.Sprintf("sensu-%v", rand.Intn(100))
		}
		if api.Host == "" {
			logger.Fatalf("Sensu API %q Host is missing", api.Name)
		}
		if api.Timeout == 0 {
			c.Sensu[i].Timeout = 10
		} else if api.Timeout >= 1000 { // backward compatibility with < 0.3.0 version
			c.Sensu[i].Timeout = api.Timeout / 1000
		}
		if api.Port == 0 {
			c.Sensu[i].Port = 4567
		}
		if api.Ssl {
			prot += "s"
		}
		c.Sensu[i].URL = fmt.Sprintf("%s://%s:%d%s", prot, api.Host, c.Sensu[i].Port, api.Path)
	}
}

func (c *Config) initGlobal() {
	if c.Dashboard != nil {
		c.Uchiwa = *c.Dashboard
	}
	if c.Uchiwa.Host == "" {
		c.Uchiwa.Host = "0.0.0.0"
	}
	if c.Uchiwa.Port == 0 {
		c.Uchiwa.Port = 3000
	}
	if c.Uchiwa.Refresh == 0 {
		c.Uchiwa.Refresh = 10
	} else if c.Uchiwa.Refresh >= 1000 { // backward compatibility with < 0.3.0 version
		c.Uchiwa.Refresh = c.Uchiwa.Refresh / 1000
	}
	if c.Uchiwa.User != "" && c.Uchiwa.Pass != "" {
		c.Uchiwa.Auth = "simple"
	}
	if c.Uchiwa.Db.Driver != "" && c.Uchiwa.Db.Scheme != "" {
		c.Uchiwa.Auth = "sql"
	}
	if c.Uchiwa.Github.Server != "" {
		c.Uchiwa.Auth = "github"
	}
}

func buildPublicConfig(c *Config) {
	p := new(Config)
	p.Uchiwa = c.Uchiwa
	p.Uchiwa.User = "*****"
	p.Uchiwa.Pass = "*****"
	p.Uchiwa.Db.Scheme = "*****"
	p.Uchiwa.Github.ClientID = "*****"
	p.Uchiwa.Github.ClientSecret = "*****"
	p.Sensu = make([]SensuConfig, len(c.Sensu))
	for i := range c.Sensu {
		p.Sensu[i] = c.Sensu[i]
		p.Sensu[i].User = "*****"
		p.Sensu[i].Pass = "*****"
	}
	PublicConfig = p
}

// LoadConfig function loads a specified configuration file and return a Config struct
func LoadConfig(path string) (*Config, error) {
	logger.Infof("Loading configuration file %s", path)
	c := new(Config)
	file, err := os.Open(path)
	if err != nil {
		if len(path) > 1 {
			return nil, fmt.Errorf("Error: could not read config file %s.", path)
		}
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return nil, fmt.Errorf("Error decoding file %s: %s", path, err)
	}

	c.initGlobal()
	c.initSensu()

	return c, nil
}
