package file

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/IliaSotnikov2005/dnsmgr/server/internal/domain"
)

type FileRepository struct {
	log  *slog.Logger
	path string
	mu   sync.RWMutex
}

func NewFileRepository(log *slog.Logger, path string) *FileRepository {
	if err := ensureFileExists(path); err != nil {
		log.Error("could not ensure file exists", "error", err)
	}

	return &FileRepository{
		log:  log,
		path: path,
	}
}

func ensureFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		return file.Close()
	}

	return nil
}

func (f *FileRepository) Get(ctx context.Context) ([]domain.DNS, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	f.log.Info("Getting DNS list from file", "path", f.path)
	file, err := os.Open(f.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			f.log.Error("failed to close the file", "error", err)
		}
	}()

	list := []domain.DNS{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "nameserver" {
			list = append(list, domain.DNS{IP: fields[1]})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return list, nil
}

func (f *FileRepository) Save(ctx context.Context, dnsList []domain.DNS) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.log.Info("Saving DNS list to file", "path", f.path)

	dir := filepath.Dir(f.path)
	tmpFile, err := os.CreateTemp(dir, ".dnsmgr.*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer func() {
		tmpFile.Close()

		os.Remove(tmpPath)
	}()

	writer := bufio.NewWriter(tmpFile)
	for _, d := range dnsList {
		if _, err := fmt.Fprintf(writer, "nameserver %s\n", d.IP); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, f.path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
