package envcontract

import (
	"fmt"

	"github.com/favxlaw/envcontract/internal/contract"
	"github.com/favxlaw/envcontract/internal/engine"
	"github.com/favxlaw/envcontract/internal/parser"
	"github.com/favxlaw/envcontract/internal/source"
)

type FieldContract = contract.FieldContract

func Validate(v any, opts ...Option) (Result, error) {
	cfg := defaultConfig()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}

	if len(cfg.sources) == 0 {
		cfg.sources = append(cfg.sources, source.SystemSource{})
	}

	contracts, err := parser.ParseStruct(v)
	if err != nil {
		return Result{}, err
	}

	env, loadWarnings, err := loadSources(cfg.sources)
	if err != nil {
		return Result{}, err
	}

	findings := engine.Run(contracts, env, engine.Options{
		CheckUnused: cfg.checkUnused,
	})

	return NewResult(findings, loadWarnings), nil
}

func loadSources(sources []Source) (map[string]string, []LoadWarning, error) {
	env := make(map[string]string)
	var warnings []LoadWarning

	for _, src := range sources {
		loaded, err := src.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("load env source: %w", err)
		}

		for key, value := range loaded.Values {
			env[key] = value
		}

		warnings = append(warnings, loaded.Warnings...)
	}

	return env, warnings, nil
}
