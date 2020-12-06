package config

// AppConfig is the configuration for the service
type AppConfig struct {
	FileStore     FileStoreConfig
	AccessControl AccessControlConfig
}

// FileStoreConfig contains the configuration for storing data in files on disk.
type FileStoreConfig struct {
	Path              string
	SystemsPerFile    int
	MemoryCacheMaxAge int
}

// AccessControlConfig defines the requirements for an incoming request to be handled.
type AccessControlConfig struct {
	RequireAccessToken bool
	AccessToken        string
}
