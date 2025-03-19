package interfaces

type ConfigProviderInterface interface {
	Get(key string) any
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetStringMap(key string) map[string]any
}
