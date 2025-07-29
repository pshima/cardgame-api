# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a Card Game API built with Go and the Gin web framework. The API provides endpoints for card game functionality.

## Tech Stack
- **Language**: Go 1.24.4
- **Framework**: Gin (github.com/gin-gonic/gin)
- **Module**: github.com/peteshima/cardgame-api

## Common Commands

### Development
- `go run main.go` - Start the development server on port 8080
- `go build` - Build the application binary
- `go mod tidy` - Clean up dependencies

### Dependencies
- `go get <package>` - Add a new dependency
- `go mod download` - Download dependencies

## Project Structure
- `main.go` - Main application entry point with server setup
- `go.mod` - Go module definition and dependencies
- `go.sum` - Dependency checksums

## API Endpoints
- `GET /hello` - Returns a JSON hello world message

## Development Notes
- Server runs on port 8080 by default
- All API responses are in JSON format
- Uses Gin's default middleware for logging and recovery