# GoDotEnvVault

![CI workflow](https://github.com/dotenv-org/godotenvvault/actions/workflows/ci.yml/badge.svg)

Extends the proven & trusted foundation of
[godotenv](https://github.com/joho/godotenv), with `.env.vault` file
support.

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

Then encrypt your environment settings by doing:

```shell
npx dotenv-vault local build
```

This will create an encrypted `.env.vault` file along with a
`.env.keys` file containing the encryption keys. Set the
`DOTENV_KEY` environment variable by copying and pasting
the key value from the `.env.keys` file:

```shell
export DOTENV_KEY="<key string from .env.keys>"
```

You can now delete your original `.env` file, and use Go like the
following to read environment settings from the encrypted `.env.vault`
file:

```go
package main

import (
    "log"
    "os"

    "github.com/dotenv-org/godotenvvault"
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

An even more convenient option is to take advantage of the autoload
package which will read in `.env.vault` on import

```go
import _ "github.com/dotenv-org/godotenvvault/autoload"
```

Note that when the `DOTENV_KEY` environment variable is set,
environment settings will *always* be loaded from the `.env.vault`
file in the project root. For development use, you can leave the
`DOTENV_KEY` environment variable unset and fall back on the
`godotenv` behaviour of loading from `.env` or a specified set of
files (see [here in the `gotodenv`
README](https://github.com/joho/godotenv#usage) for the details).

If you don't want `godotenvvault` to modify your program's environment
directly, you can just load and decrypt the `.env.vault` file and get
the result as a map by doing:

```go
var myEnv map[string]string
myEnv, err := godotenvvault.Read()

s3Bucket := myEnv["S3_BUCKET"]
```

## FAQ

#### What happens if `DOTENV_KEY` is not set?

Dotenv Vault gracefully falls back to
[godotenv](https://github.com/joho/godotenv) when `DOTENV_KEY` is not
set. This is the default for development so that you can focus on
editing your `.env` file and save the `build` command until you are
ready to deploy those environment variables changes.

#### Should I commit my `.env` file?

No. We **strongly** recommend against committing your `.env` file to
version control. It should only include environment-specific values
such as database passwords or API keys. Your production database
should have a different password than your development database.

#### Should I commit my `.env.vault` file?

Yes. It is safe and recommended to do so. It contains your encrypted
envs, and your vault identifier.

#### Can I share the `DOTENV_KEY`?

No. It is the key that unlocks your encrypted environment variables.
Be very careful who you share this key with. Do not let it leak.
