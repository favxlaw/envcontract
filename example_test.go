package envcontract_test

import (
	"fmt"

	"github.com/favxlaw/envcontract"
)

func ExampleValidate() {
	type Config struct {
		Host string `env:"HOST,default=localhost"`
		Port int    `env:"PORT,required"`
	}

	src := staticSource{
		values: map[string]string{
			"HOST": "localhost",
			"PORT": "8080",
		},
	}

	result, err := envcontract.Validate(&Config{}, envcontract.WithSource(src))
	if err != nil {
		panic(err)
	}

	fmt.Println(result.HasErrors())

	// Output:
	// false
}

type staticSource struct {
	values map[string]string
}

func (s staticSource) Load() (envcontract.LoadResult, error) {
	return envcontract.LoadResult{
		Values: s.values,
	}, nil
}
