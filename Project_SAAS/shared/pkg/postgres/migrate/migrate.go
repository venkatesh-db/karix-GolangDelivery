package migrate

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run executes the SQL files in dir in lexical order against the provided pool.
func Run(ctx context.Context, pool *pgxpool.Pool, filesystem fs.FS, dir string) error {
	entries, err := fs.ReadDir(filesystem, dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		fullPath := path.Join(dir, entry.Name())
		content, err := fs.ReadFile(filesystem, fullPath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("apply migration %s: %w", entry.Name(), err)
		}
	}
	return nil
}
