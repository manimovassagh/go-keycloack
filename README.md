
# Keycloak Secured Go Application

This project demonstrates a simple Go application secured with Keycloak. It uses Echo as the web framework and Gocloak as the Keycloak client.

## Features

- Public and secure endpoints
- JWT token validation using Keycloak
- Modular project structure

## Table of Contents

- [Keycloak Secured Go Application](#keycloak-secured-go-application)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Getting Started](#getting-started)
    - [Clone the Repository](#clone-the-repository)
    - [Set Up Environment Variables](#set-up-environment-variables)
    - [Run with Docker Compose](#run-with-docker-compose)
  - [Project Structure](#project-structure)
  - [Usage](#usage)
  - [Endpoints](#endpoints)
  - [Contributing](#contributing)
  - [License](#license)

## Prerequisites

- Docker
- Docker Compose

## Getting Started

### Clone the Repository

```sh
git clone https://github.com/yourusername/your_project.git
cd your_project
```

### Set Up Environment Variables

Create a `.env` file in the project root directory and set the necessary environment variables:
Obviously you should use your own variables. the followings are only samples :)

```sh
KEYCLOAK_URL=http://localhost:8080
REALM=mani
CLIENT_ID=go
CLIENT_SECRET=aZbUGUvcXxETMznmtsFB1Mm8TEWK4Hpz
```

### Run with Docker Compose

Build and run the application using Docker Compose:

```sh
docker compose up --build
```

This will start the Keycloak server and the Go application.

## Project Structure

```
your_project/
|-- config/
|   |-- config.go
|-- handlers/
|   |-- handlers.go
|-- auth/
|   |-- auth.go
|-- main.go
|-- go.mod
|-- go.sum
|-- Dockerfile
|-- docker-compose.yml
```

- **config/**: Contains configuration-related code.
- **handlers/**: Contains HTTP handler functions.
- **auth/**: Contains authentication middleware and related functions.
- **main.go**: Entry point of the application.
- **Dockerfile**: Dockerfile for building the Go application.
- **docker-compose.yml**: Docker Compose configuration file.

## Usage

### Access Public Endpoint

You can access the public endpoint without authentication:

```sh
curl http://localhost:4000/public
```

### Access Secure Endpoint

To access the secure endpoint, you need to obtain a JWT token from Keycloak and include it in the request header.

```sh
curl -H "Authorization: Bearer <your_access_token>" http://localhost:4000/secure
```

## Endpoints

- **GET /public**: Public endpoint accessible without authentication.
- **GET /secure**: Secure endpoint accessible with a valid JWT token.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any bugs, improvements, or new features.

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
