//go:build !release

package internal

// This file ensures that environment variables from .env files are
// automatically loaded on application start.
// As this file is exluded from the release build, this behaviour will only
// be disabled in Docker Images and builds using the release tag

import _ "github.com/joho/godotenv/autoload"
