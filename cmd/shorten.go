package main

import (
	"micro-url/internal/shortener"
	"micro-url/internal/shortener/repo/etcd"
	"micro-url/internal/shortener/service"
	"micro-url/internal/shortener/transport"

	"github.com/caarlos0/env/v6"
)

func main() {
	cfg := shortener.Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	ob, err := service.NewKnuthHashidsObfuscator(cfg.Crypto)
	if err != nil {
		panic(err)
	}

	r, err := etcd.NewEtcdRepo(cfg.EtcdClient)
	if err != nil {
		panic(err)
	}

	svc := service.NewURLShortener(r, ob)

	errCh := make(chan error)
	go transport.Start(cfg, svc, errCh)
	panic(<-errCh)
}
