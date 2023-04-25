# GoDotEnvVault

![CI workflow](https://github.com/dotenv-org/godotenvvault/actions/workflows/ci.yml/badge.svg)

Extends the proven & trusted foundation of
[godotenv](https://github.com/joho/godotenv), with `.env.vault` file
support.

* [🌱 Install](#-install)
* [🏗️ Usage (.env)](#-usage)
* [🚀 Deploying (.env.vault) 🆕](#-deploying)
* [🌴 Multiple Environments](#-manage-multiple-environments)
* [❓ FAQ](#-faq)
* [⏱️ Changelog](./CHANGELOG.md)

## 🌱 Install

```shell
go get github.com/dotenv-org/godotenvvault
```

## 🏗️ Usage

Add your application configuration to your `.env` file in the root of your project:

```shell
# .env
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
```

As early as possible in your application, import and configure godotenvvault:

```
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

That's it! `os.Getenv` has the keys and values you defined in your .env file. Continue using it this way in development. It works just like [godotenv](https://github.com/joho/godotenv). In the next section, we'll unlock the recommended deployment mechanism using `.env.vault`.

## 🚀 Deploying

Encrypt your environment settings by doing:

```shell
npx dotenv-vault local build
```

This will create an encrypted `.env.vault` file along with a
`.env.keys` file containing the encryption keys. Set the
`DOTENV_KEY` environment variable by copying and pasting
the key value from the `.env.keys` file onto your server
or cloud provider. For example in heroku:

```shell
heroku config:set DOTENV_KEY=<key string from .env.keys>
```

Commit your .env.vault file safely to code and deploy. Your .env.vault fill be decrypted on boot, its environment variables injected, and your app work as expected.

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

## 🌴 Manage Multiple Environments

Create a `.env.production` file in the root of your project and put your production values there.

```shell
# .env.production
S3_BUCKET="PRODUCTION_S3BUCKET"
SECRET_KEY="PRODUCTION_SECRETKEYGOESHERE"
```

Rebuild your `.env.vault` file.

```shell
npx dotenv-vault local build
```

View your `.env.keys` file. There is a production `DOTENV_KEY` that coincides with the additional `DOTENV_VAULT_PRODUCTION` cipher in your `.env.vault` file.

Set the production `DOTENV_KEY` on your server, recommit your `.env.vault` file to code, and deploy. That's it! Your .env.vault fill be decrypted on boot, its production environment variables injected, and your app work as expected.

Want to additionally backup your .env files, maintain access controls, change history, and more? Check out the [vault managed guide to multiple environments](https://www.dotenv.org/docs/languages/go#-manage-multiple-environments).

## ❓ FAQ

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
