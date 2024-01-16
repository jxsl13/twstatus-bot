package db

import "io/fs"

type Options struct {
	ssl          bool
	migrationsFs fs.FS
	versionTable string
}

func (o *Options) SSL() string {
	if o.ssl {
		return "disable"
	}
	return "enable"
}

type Option func(*Options)

func WithSSL(ssl bool) Option {
	return func(o *Options) {
		o.ssl = ssl
	}
}

func WithMigrationsFs(migrationsFs fs.FS) Option {
	return func(o *Options) {
		o.migrationsFs = migrationsFs
	}
}

func WithVersionTable(versionTable string) Option {
	return func(o *Options) {
		o.versionTable = versionTable
	}
}
