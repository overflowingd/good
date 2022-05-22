package dotenv

import "github.com/joho/godotenv"

func Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}
