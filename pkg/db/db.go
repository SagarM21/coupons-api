package db

import (
	"fmt"
	"strings"
	"time"

	"e-commerce/pkg/config"

	"github.com/gocql/gocql"
)

type DB struct {
	Session *gocql.Session
}

var consistencyMap = map[string]gocql.Consistency{
	"any":          gocql.Any,
	"one":          gocql.One,
	"two":          gocql.Two,
	"three":        gocql.Three,
	"quorum":       gocql.Quorum,
	"all":          gocql.All,
	"localquorum":  gocql.LocalQuorum,
	"eachquorum":   gocql.EachQuorum,
	"localone":     gocql.LocalOne,
}

func Init() (*DB, error) {
	cfg := config.App.Database

	cluster := gocql.NewCluster(cfg.Host)
	cluster.Port = cfg.Port
	cluster.Keyspace = cfg.Keyspace
	cluster.ConnectTimeout = 30 * time.Second
	cluster.Timeout = 30 * time.Second

	if c, ok := consistencyMap[strings.ToLower(cfg.Consistency)]; ok {
		cluster.Consistency = c
	} else {
		cluster.Consistency = gocql.One
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to scylla: %w", err)
	}

	return &DB{Session: session}, nil
}

func (d *DB) Close() {
	d.Session.Close()
}

func CreateKeyspace(cfg config.Database) error {
	cluster := gocql.NewCluster(cfg.Host)
	cluster.Port = cfg.Port
	cluster.Consistency = gocql.All
	cluster.ConnectTimeout = 30 * time.Second
	cluster.Timeout = 30 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to connect to scylla for keyspace creation: %w", err)
	}
	defer session.Close()

	query := fmt.Sprintf(
		`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`,
		cfg.Keyspace,
	)
	return session.Query(query).Exec()
}

func (d *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS coupons (
			id TEXT PRIMARY KEY,
			type TEXT,
			details TEXT,
			created_at TIMESTAMP,
			expires_at TIMESTAMP
		)`,
	}

	for _, q := range queries {
		if err := d.Session.Query(q).Exec(); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}
	return nil
}
