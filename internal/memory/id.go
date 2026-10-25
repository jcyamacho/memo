package memory

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	maxNameLength = 48
	maxIDLength   = 64
	suffixBytes   = 3
)

var namedIDPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type ID struct {
	value string
}

func NewID() (ID, error) {
	value, err := uuid.NewRandom()
	if err != nil {
		return ID{}, err
	}
	return ID{value: value.String()}, nil
}

func NewNamedID(name string) (ID, error) {
	slug := normalizeName(name)
	if slug == "" {
		return ID{}, ErrInvalidInput
	}

	suffix := make([]byte, suffixBytes)
	rand.Read(suffix)
	return ID{value: slug + "-" + hex.EncodeToString(suffix)}, nil
}

func ParseID(value string) (ID, error) {
	if parsed, err := uuid.Parse(value); err == nil {
		if parsed == uuid.Nil {
			return ID{}, ErrInvalidInput
		}
		return ID{value: parsed.String()}, nil
	}

	if len(value) > maxIDLength || !namedIDPattern.MatchString(value) {
		return ID{}, ErrInvalidInput
	}
	return ID{value: value}, nil
}

func (id ID) String() string {
	return id.value
}

func normalizeName(name string) string {
	var slug strings.Builder
	pendingDash := false
	for _, r := range strings.ToLower(name) {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			pendingDash = true
			continue
		}
		if pendingDash && slug.Len() > 0 {
			slug.WriteByte('-')
		}
		pendingDash = false
		slug.WriteRune(r)
	}

	normalized := slug.String()
	if len(normalized) > maxNameLength {
		normalized = strings.TrimRight(normalized[:maxNameLength], "-")
	}
	return normalized
}
