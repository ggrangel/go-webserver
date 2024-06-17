##  Overview

This is a Go web server project aimed at learning and implementing best practices for building robust APIs. It follows an onion architecture for clear separation of concerns.

## Key Features & Learning Objectives

- **Architecture**: Onion architecture to promote maintainability and testability.
- **Authentication**/Authorization: JWTs and refresh tokens for secure access control.
- **Data Management**: CRUD operations on an in-memory (JSON) database. Password hashing for security.
- **API Endpoints**: RESTful design with metrics, health checks, and webhooks.
- **Developer Experience**: Automated server restarts using Air.

##  Getting Started

### Prerequistes

1. Go 1.22 or higher
2. Air (for hot-reloading)

###  Installation

> 1. Clone the go-webserver repository and install the dependencies:
> ```console
> $ go build 
> ```

Set the environment variables in the `.env` file. The following variables are required:
1. `JWT_SECRET`: Secret key for JWT token generation.
2. `POLKA_KEY`: Fake API key to serve as an authorization mechanism.

###  Usage

> Start the server on port 8080 with the following command
> ```console
> $ air -c .air.toml
> ```

Use your favorite HTTP client to interact with the API. 
Check `main.go` to see the available routes and their handles for the expected request payloads.

---
