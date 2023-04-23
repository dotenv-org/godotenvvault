package autoload

import "github.com/dotenv-org/godotenvvault"

// You can just read the .env file on import just by doing
//		import _ "github.com/dotenv-org/godotenvvault/autoload"

func init() {
	godotenvvault.Load()
}
