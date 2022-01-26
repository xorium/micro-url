package service

import (
	"micro-url/internal/shortener"
	"strconv"

	"github.com/pjebs/optimus-go"
	"github.com/speps/go-hashids/v2"
)

type IDObfuscator interface {
	// Obfuscate perform padding for little ID values and simple obfuscating on
	// the numbers transform level. After that special encoding function applies
	// to get short string ID.
	Obfuscate(id uint64) string
}

// KnuthHashidsObfuscator based on basic counter value by convert to modular representation
// and hashids (http://www.hashids.org/).
type KnuthHashidsObfuscator struct {
	engine optimus.Optimus
	hasher *hashids.HashID
}

func NewKnuthHashidsObfuscator(cfg shortener.CfgCrypto) (*KnuthHashidsObfuscator, error) {
	modInverse := cfg.ModInverse
	prime := cfg.Prime
	random := cfg.Random
	if modInverse == 0 {
		modInverse = optimus.ModInverse(prime)
	}

	hd := hashids.NewData()
	hd.Salt = strconv.FormatUint(prime^random, 16) //nolint:gomnd
	hd.MinLength = 1
	hasher, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}

	return &KnuthHashidsObfuscator{
		engine: optimus.New(prime, modInverse, random),
		hasher: hasher,
	}, nil
}

func (o *KnuthHashidsObfuscator) Obfuscate(id uint64) string {
	obfuscated := o.engine.Encode(id)
	encoded, _ := o.hasher.Encode([]int{int(obfuscated)})
	return encoded
}
