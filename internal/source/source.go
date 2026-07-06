package source

type Source interface {
	Load() (LoadResult, error)
}

type LoadResult struct {
	Values   map[string]string
	Warnings []LoadWarning
}

type LoadWarning struct {
	Line    int
	Message string
}
