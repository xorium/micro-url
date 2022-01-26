package service

import (
	"context"
	"fmt"
	"micro-url/internal/shortener/repo"
	"strings"

	"github.com/PuerkitoBio/purell"
)

const serviceName = "url-shortener"

type URLShortener struct {
	repo repo.URLRepo
	ob   IDObfuscator
}

func NewURLShortener(r repo.URLRepo, o IDObfuscator) *URLShortener {
	return &URLShortener{
		repo: r,
		ob:   o,
	}
}

func (sh *URLShortener) ShortenURL(ctx context.Context, u string) (shortenID string, err error) {
	normalized, err := purell.NormalizeURLString(
		u, purell.FlagLowercaseScheme|purell.FlagLowercaseHost|purell.FlagUppercaseEscapes|
			purell.FlagsUsuallySafeGreedy,
	)
	if err != nil {
		return "", fmt.Errorf("can't normalize URL: %w", err)
	}
	if !strings.HasPrefix(normalized, "http") {
		return "", fmt.Errorf("invalid URL schema: must be a HTTP URL")
	}

	numValue, err := sh.repo.IncrementLatestCounterValue(ctx, serviceName, 1)
	if err != nil {
		return "", fmt.Errorf("can't prepare all data for the shortenning: %w", err)
	}
	// Obfuscating basic counter value by convert to modular representation
	// (D.Knuth, Vol. 2, Chapter 4.3.2) and apply hashids library.
	return sh.ob.Obfuscate(numValue), nil
}
