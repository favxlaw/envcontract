package envcontract

import (
	"encoding/json"
	"fmt"

	"github.com/favxlaw/envcontract/internal/parser"
)

// ExportSchema derives v's field contract and encodes it as indented JSON.
func ExportSchema(v any) ([]byte, error) {
	contracts, err := parser.ParseStruct(v)
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(contracts, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("envcontract: marshal schema: %w", err)
	}

	return data, nil
}

// ParseSchema decodes JSON produced by ExportSchema back into field contracts.
func ParseSchema(data []byte) ([]FieldContract, error) {
	var contracts []FieldContract

	if err := json.Unmarshal(data, &contracts); err != nil {
		return nil, fmt.Errorf("envcontract: parse schema: %w", err)
	}

	return contracts, nil
}
