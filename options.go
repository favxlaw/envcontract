package envcontract

import "github.com/favxlaw/envcontract/internal/source"

type Source interface {
	Load() (LoadResult, error)
}

type LoadResult = source.LoadResult

type LoadWarning = source.LoadWarning

type Option func(*config)

type config struct {
	sources     []Source
	checkUnused bool
}

func defaultConfig() config {
	return config{}
}

func WithFile(path string) Option {
	return func(cfg *config) {
		cfg.sources = append(cfg.sources, source.FileSource{Path: path})
	}
}

func WithSystemEnv() Option {
	return func(cfg *config) {
		cfg.sources = append(cfg.sources, source.SystemSource{})
	}
}

func WithSource(src Source) Option {
	return func(cfg *config) {
		if src == nil {
			return
		}

		cfg.sources = append(cfg.sources, src)
	}
}

func WithUnusedCheck() Option {
	return func(cfg *config) {
		cfg.checkUnused = true
	}
}
