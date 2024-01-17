package db

import (
	"io/fs"
	"time"
)

type Options struct {
	sslmode      SSLMode
	migrationsFs fs.FS
	versionTable string
	connTimeout  time.Duration
}

type Option func(*Options) error

func WithSSL(ssl SSLMode) Option {
	return func(o *Options) error {
		var sslmode SSLMode
		err := sslmode.UnmarshalText([]byte(ssl))
		if err != nil {
			return err
		}
		o.sslmode = sslmode
		return nil
	}
}

func WithMigrationsFs(migrationsFs fs.FS) Option {
	return func(o *Options) error {
		o.migrationsFs = migrationsFs
		return nil
	}
}

func WithVersionTable(versionTable string) Option {
	return func(o *Options) error {
		o.versionTable = versionTable
		return nil
	}
}

func WithConnTimeout(connTimeout time.Duration) Option {
	return func(o *Options) error {
		if connTimeout < time.Second*1 {
			o.connTimeout = time.Second * 1
		} else {
			o.connTimeout = connTimeout
		}
		return nil
	}
}
