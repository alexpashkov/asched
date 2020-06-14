package env

import "os"

func PORT() string {
	return os.Getenv("PORT")
}

func MONGODB_URI() string {
	return os.Getenv("MONGODB_URI") + "?retryWrites=false"
}

func MONGODB_DB_NAME() string {
	return os.Getenv("MONGODB_DB_NAME")
}
