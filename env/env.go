package env

import "os"

func DbType() string {
	return os.Getenv("DB_TYPE")
}

func DbName() string {
	return os.Getenv("DB_NAME")
}

func DbDsn() string {
	return os.Getenv("DB_DSN")
}

func Port() string {
	return os.Getenv("PORT")
}

func JwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func SessionSecret() string {
	return os.Getenv("SESSION_SECRET")
}

func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func SecretFlag() string {
	return os.Getenv("SECRET_FLAG")
}

func ProblemRoot() string {
	return os.Getenv("PROBLEM_ROOT")
}
