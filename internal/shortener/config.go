package shortener

import "time"

type Config struct {
	EtcdClient CfgEtcdClient `envPrefix:"ETCD_"`
	HTTTPAddr  string        `env:"HTTP" envDefault:":8888"`
	Crypto     CfgCrypto     `envPrefix:"CRYPTO_"`
}

type CfgEtcdClient struct {
	Peers   []string      `env:"PEERS" envSeparator:","`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"2s"`
}

type CfgCrypto struct {
	Prime      uint64 `env:"PRIME,notEmpty,unset"`
	ModInverse uint64 `env:"MOD_INVERSE,unset" envDefault:"0"`
	Random     uint64 `env:"RANDOM,notEmpty,unset"`
	Alphabet   string `env:"ALPHA,notEmpty"`
}
