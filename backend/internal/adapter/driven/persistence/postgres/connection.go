package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Url string
}

// MustConnect establishes a connection or panics
func MustConnect(cfg *Config) *gorm.DB {
	db, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return db
}

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Url), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		NowFunc:     time.Now,
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(50)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(50)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// stats := sqlDB.Stats()
	// fmt.Printf("Conexões abertas: %d\n", stats.OpenConnections)
	// fmt.Printf("Conexões em uso: %d\n", stats.InUse)
	// fmt.Printf("Conexões inativas: %d\n", stats.Idle)
	// fmt.Printf("Vezes esperando por conexões: %d\n", stats.WaitCount)
	// fmt.Printf("Tempo esperando por conexões: %s\n", stats.WaitDuration)
	// fmt.Printf("Conexões fechadas por inatividade: %d\n", stats.MaxIdleClosed)
	// fmt.Printf("Conexões fechadas por tempo de vida: %d\n", stats.MaxLifetimeClosed)

	return db, nil
}

// OpenConnections: O número total de conexões abertas atualmente, incluindo tanto as inativas quanto as ativas.
// InUse: O número de conexões atualmente em uso, ou seja, conexões que estão sendo ativamente utilizadas pela aplicação.
// Idle: O número de conexões inativas (disponíveis para reutilização no pool).
// WaitCount: O número de vezes que a aplicação precisou esperar por uma nova conexão porque todas as conexões abertas estavam em uso. Esse número pode indicar que o limite de conexões abertas foi atingido.
// WaitDuration: O tempo total que a aplicação passou esperando por uma conexão disponível. Um valor alto pode indicar que o sistema está atingindo o limite de conexões com frequência.
// MaxIdleClosed: O número de conexões que foram fechadas automaticamente porque estavam ociosas por muito tempo e ultrapassaram o limite definido em SetConnMaxIdleTime.
// MaxLifetimeClosed: O número de conexões que foram fechadas porque atingiram o limite máximo de tempo de vida definido em SetConnMaxLifetime.
