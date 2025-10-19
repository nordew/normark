package config

type Config struct {
	Server    Server
	Postgres  Postgres
	JWT       JWT
	CORS      CORS
	RateLimit RateLimit
}

type Server struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Database string `env:"POSTGRES_DB" envDefault:"postgres"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
}

type JWT struct {
	Secret             string `env:"JWT_SECRET,required"`
	AccessTokenExpiry  int    `env:"JWT_ACCESS_TOKEN_EXPIRY" envDefault:"15"`
	RefreshTokenExpiry int    `env:"JWT_REFRESH_TOKEN_EXPIRY" envDefault:"10080"`
}

type CORS struct {
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000"`
	AllowMethods     []string `env:"CORS_ALLOW_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders     []string `env:"CORS_ALLOW_HEADERS" envSeparator:"," envDefault:"Origin,Content-Type,Authorization"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
	MaxAge           int      `env:"CORS_MAX_AGE" envDefault:"43200"`
}

type RateLimit struct {
	RequestsPerSecond int `env:"RATE_LIMIT_RPS" envDefault:"10"`
	Burst             int `env:"RATE_LIMIT_BURST" envDefault:"20"`
}
