package source

type Source interface {
	Load() (map[string]string, error)
}
