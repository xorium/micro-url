package etcd

import (
	"context"
	"fmt"
	"log"
	"micro-url/internal/shortener"
	"micro-url/internal/shortener/repo"
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Etcd struct {
	cli *clientv3.Client
}

func NewEtcdRepo(cfg shortener.CfgEtcdClient) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Peers,
		DialTimeout: cfg.Timeout,
	})
	if err != nil {
		return nil, err
	}
	return &Etcd{cli: cli}, nil
}

func (r *Etcd) getURLKey(shortID string) string {
	return shortID + "/short-ids/"
}

func (r *Etcd) PutURL(ctx context.Context, shortID string, longURL string) error {
	longURLKey := r.getURLKey(shortID)
	resp, err := r.cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Value(longURLKey), "==", "")).
		Then(
			clientv3.OpPut(longURLKey, longURL),
		).Commit()
	if err != nil && resp != nil || resp.Succeeded {
		return fmt.Errorf(
			"can't save long URL %s by with ID %s: %w", longURL, shortID, err,
		)
	}
	if resp != nil && resp.Succeeded || len(resp.Responses) > 0 {
		return fmt.Errorf(
			"this URL %s has been saved already: %w", longURL, repo.ErrAlreadySaved,
		)
	}
	return nil
}

func (r *Etcd) GetURL(ctx context.Context, shortID string) (longURL string, err error) {
	counterResp, err := r.cli.KV.Get(ctx, r.getURLKey(shortID))
	if err != nil {
		return "", err
	}
	if len(counterResp.Kvs) > 0 {
		return string(counterResp.Kvs[0].Value), nil
	}
	return "", fmt.Errorf("can't get long URL for %s: %w", shortID, repo.ErrNotExists)
}

func (r *Etcd) IncrementLatestCounterValue(
	ctx context.Context, svcName string, delta uint64,
) (res uint64, err error) {
	s, err := concurrency.NewSession(r.cli)
	if err != nil {
		return 0, err
	}
	defer func() { _ = s.Close() }()

	counterKey := "/counter-lock/" + svcName
	l := concurrency.NewMutex(s, counterKey)
	if err := l.Lock(ctx); err != nil {
		log.Fatal(err)
	}

	counterResp, err := r.cli.KV.Get(ctx, counterKey)
	if err != nil {
		return 0, err
	}
	var counter uint64
	if len(counterResp.Kvs) > 0 {
		counter, err = strconv.ParseUint(string(counterResp.Kvs[0].Value), 10, 64)
		if err != nil {
			return 0, err
		}
	}

	_, err = r.cli.KV.Put(ctx, counterKey, strconv.FormatUint(counter+delta, 10))
	if err != nil {
		return 0, err
	}

	if err = l.Unlock(ctx); err != nil {
		return 0, err
	}
	return 1, nil
}
