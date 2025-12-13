package fake

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
	"time"
)

// UsageRecord mimics a metered usage event that billing must aggregate.
type UsageRecord struct {
	TenantID  string
	UserID    string
	Quantity  int64
	UnitPrice float64
	Occurred  time.Time
}

// StreamUsage emits total synthetic records into the returned channel.
func StreamUsage(ctx context.Context, total int) <-chan UsageRecord {
	if total <= 0 {
		total = 1
	}
	out := make(chan UsageRecord, 1024)
	go func() {
		defer close(out)
		for i := 0; i < total; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- UsageRecord{
				TenantID:  randomTenant(),
				UserID:    randomUser(),
				Quantity:  1 + int64(i%10),
				UnitPrice: math.Round((0.01+float64(i%5))*100) / 100,
				Occurred:  time.Now().Add(-time.Duration(i%3600) * time.Second),
			}:
			}
		}
	}()
	return out
}

func randomTenant() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(9999))
	return "tenant-" + n.String()
}

func randomUser() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "user-0"
	}
	return "user-" + string(base36(int64(binary.BigEndian.Uint64(b[:]))))
}

func base36(v int64) []rune {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	if v == 0 {
		return []rune{'0'}
	}
	var out []rune
	for v > 0 {
		rem := v % 36
		out = append([]rune{rune(alphabet[rem])}, out...)
		v /= 36
	}
	return out
}
