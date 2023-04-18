# GoDotEnvVault

Extends the proven & trusted foundation of [godotenv](https://github.com/joho/godotenv), with `.env.vault` file support.

## Installation

## Usage

## FAQ

#### What happens if `DOTENV_KEY` is not set?

Dotenv Vault gracefully falls back to [godotenv](https://github.com/joho/godotenv) when `DOTENV_KEY` is not set. This is the default for development so that you can focus on editing your `.env` file and save the `build` command until you are ready to deploy those environment variables changes.

#### Should I commit my `.env` file?

No. We **strongly** recommend against committing your `.env` file to version control. It should only include environment-specific values such as database passwords or API keys. Your production database should have a different password than your development database.

#### Should I commit my `.env.vault` file?

Yes. It is safe and recommended to do so. It contains your encrypted envs, and your vault identifier.

#### Can I share the `DOTENV_KEY`?

No. It is the key that unlocks your encrypted environment variables. Be very careful who you share this key with. Do not let it leak.
