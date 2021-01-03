package storage

// Config struct for database connection
type Config struct {
	Host     string
	Username string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}
