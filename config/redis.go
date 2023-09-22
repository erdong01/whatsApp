package config

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
}

type RedisCluster struct {
	Password string   `mapstructure:"password" json:"password" yaml:"password"` // 密码
	Addrs    []string `mapstructure:"addrs" json:"addrs" yaml:"addrs"`          // 服务器地址:端口
}
