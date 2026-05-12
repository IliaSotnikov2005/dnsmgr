package file

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/IliaSotnikov2005/dnsmgr/server/internal/domain"
)

type FileRepository struct {
	log  *slog.Logger
	path string
}

func NewFileRepository(log *slog.Logger, path string) *FileRepository {
	return &FileRepository{
		log:  log,
		path: path,
	}
}

func (f *FileRepository) Get(ctx context.Context) ([]domain.DNS, error) {
	f.log.Info("Getting DNS list from file", "path", f.path)
	file, err := os.Open(f.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var list []domain.DNS
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
	f.log.Info("Saving DNS list to file", "path", f.path)
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, d := range dnsList {
		if _, err := fmt.Fprintf(writer, "nameserver %s\n", d.IP); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}

	return nil
}
