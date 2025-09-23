package db

import "gorm.io/gorm"

// Options represents configuration options for database operations
type Options struct {
	Tx      *gorm.DB
	Preload bool
}

// Option is a function type that modifies Options
type Option func(*Options)

// WithTx returns an Option that sets the database transaction
func WithTx(tx *gorm.DB) Option {
	return func(o *Options) {
		o.Tx = tx
	}
}

// WithPreload returns an Option that enables preloading of related data
func WithPreload() Option {
	return func(o *Options) {
		o.Preload = true
	}
}
