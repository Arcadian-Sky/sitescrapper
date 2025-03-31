package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Arcadian-Sky/scrapper/internal/models"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseStorage struct {
	conn driver.Conn
}

func NewClickHouseStorage(host, port, database, username, password string) (*ClickHouseStorage, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к ClickHouse: %v", err)
	}

	// Проверяем подключение
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к ClickHouse: %v", err)
	}

	return &ClickHouseStorage{conn: conn}, nil
}

func (s *ClickHouseStorage) SaveProduct(ctx context.Context, product models.Product) error {
	query := `
		INSERT INTO products (name, description, price, url, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	return s.conn.Exec(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.URL,
		time.Now(),
	)
}

func (s *ClickHouseStorage) Close() error {
	return s.conn.Close()
}
