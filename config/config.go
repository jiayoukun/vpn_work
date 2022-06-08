package config

import (
	"context"
	"errors"
	"github.com/vishvananda/netlink"
)

type Config struct {
	Server struct{
		Name	string	`yaml:"name"`
		Address	string	`yaml:"address"`
		MaxWorkers	int	`yaml:"maxworkers"`
		Keepalive	int	`yaml:"keepalive"`
		Insecure	bool `yaml:"insecure"`
		Mtu	int`yaml:"mtu"`
	}`yaml:"server"`

	Crypto struct{
		Type string	`yaml:"type"`
		Key string	`yaml:"key"`
	}`yaml:"crypto"`

	Nodes	[]struct{
		Node `yaml:"node"`
	}`yaml:"nodes"`

	source source
}

type Node struct {
	Name string `yaml:"name"`
	Address          string   `yaml:"address"`
	PrivateAddresses []string `yaml:"privateAddresses"`
	PrivateSubnets   []string `yaml:"privateSubnets"`
}

type source interface {
	load() (*Config, error)
}

func New() *Config {
	return &Config{}
}

func (c *Config)FromFile(cfile string) *Config {
	c.source = &file{
		path: "config.yaml",
		cfile: cfile,
	}
	return c
}

func (c *Config) Load() error {
	cfg, err := c.source.load()
	if err != nil {
		return err
	}
	source := c.source
	*c = *cfg
	c.source = source

	setDefaultConfig(c)

	return nil
}

func (c *Config) Wacher(ctx context.Context, extNotify chan struct{})  {
	notify := make(chan struct{}, 1)

	go func(n chan struct{}) {
		for {
			select {
			case <- notify:
			case <- ctx.Done():
				return
			}
			c.Load()
			extNotify <- struct{}{}

		}
	}(extNotify)
}

// Search returns current node config
func (c Config)Search() (Node, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return Node{}, err
	}
	var ipList []string
	for _, link := range links {
		addres, err := netlink.AddrList(link, netlink.FAMILY_ALL)//FAMILY_ALL  linux常量
		if err != nil {
			return Node{}, err
		}
		for	_, addr := range addres {
			ipList = append(ipList, addr.IP.String())
		}
	}

	for _, nodes := range c.Nodes {
		for _, ip := range ipList {
			if nodes.Node.Address == ip {
				return nodes.Node, nil
			}
		}
	}
	return Node{}, errors.New("Search error: can not find node")
}

func setDefaultConfig(c *Config) {
	// set defaults
	if c.Server.Address == "" {
		c.Server.Address = ":8085"
	}

	if c.Server.MaxWorkers == 0 {
		c.Server.MaxWorkers = 10
	}

	if c.Server.Keepalive == 0 {
		c.Server.Keepalive = 10
	}

	if c.Server.Mtu == 0 {
		c.Server.Mtu = 1300
	}
}
