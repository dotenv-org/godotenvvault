# GoDotEnvVault

Extends the proven & trusted foundation of [godotenv](https://github.com/joho/godotenv), with `.env.vault` file support.

## Installation

As a library

```shell
go get github.com/dotenv-org/godotenvvault
```

## Usage

Add your application configuration to your `.env` file in the root of your project:

```shell
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
```

Then in your Go app you can do something like

```go
package main

import (
    "log"
    "os"

    "github.com/joho/godotenvvault"
)

func main() {
  err := godotenvvault.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  s3Bucket := os.Getenv("S3_BUCKET")
  secretKey := os.Getenv("SECRET_KEY")

  // now do something with s3 or whatever
}
```

If you're even lazier than that, you can just take advantage of the autoload package which will read in `.env` on import

```go
import _ "github.com/dotenv-org/godotenvvault/autoload"
```

While `.env` in the project root is the default, you don't have to be constrained, both examples below are 100% legit

```go
godotenvvault.Load("somerandomfile")
godotenvvault.Load("filenumberone.env", "filenumbertwo.env")
```

If you want to be really fancy with your env file you can do comments and exports (below is a valid env file)

```shell
# I am a comment and that is OK
SOME_VAR=someval
FOO=BAR # comments at line end are OK too
export BAR=BAZ
```

Or finally you can do YAML(ish) style

```yaml
FOO: bar
BAR: baz
```

as a final aside, if you don't want godotenvvault munging your env you can just get a map back instead

```go
var myEnv map[string]string
myEnv, err := godotenvvault.Read()

s3Bucket := myEnv["S3_BUCKET"]
```

... or from an `io.Reader` instead of a local file

```go
reader := getRemoteFile()
myEnv, err := godotenvvault.Parse(reader)
```

... or from a `string` if you so desire

```go
content := getRemoteFileContent()
myEnv, err := godotenvvault.Unmarshal(content)
```

### Precedence & Conventions

Existing envs take precedence of envs that are loaded later.

The [convention](https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use)
for managing multiple environments (i.e. development, test, production)
is to create an env named `{YOURAPP}_ENV` and load envs in this order:

```go
env := os.Getenv("FOO_ENV")
if "" == env {
  env = "development"
}

godotenvvault.Load(".env." + env + ".local")
if "test" != env {
  godotenvvault.Load(".env.local")
}
godotenvvault.Load(".env." + env)
godotenvvault.Load() // The Original .env
```

If you need to, you can also use `godotenvvault.Overload()` to defy this convention
and overwrite existing envs instead of only supplanting them. Use with caution.

## FAQ

#### What happens if `DOTENV_KEY` is not set?

Dotenv Vault gracefully falls back to [godotenv](https://github.com/joho/godotenv) when `DOTENV_KEY` is not set. This is the default for development so that you can focus on editing your `.env` file and save the `build` command until you are ready to deploy those environment variables changes.

#### Should I commit my `.env` file?

No. We **strongly** recommend against committing your `.env` file to version control. It should only include environment-specific values such as database passwords or API keys. Your production database should have a different password than your development database.

#### Should I commit my `.env.vault` file?

Yes. It is safe and recommended to do so. It contains your encrypted envs, and your vault identifier.

#### Can I share the `DOTENV_KEY`?

No. It is the key that unlocks your encrypted environment variables. Be very careful who you share this key with. Do not let it leak.
