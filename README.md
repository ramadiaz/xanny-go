<a id="readme-top"></a>

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue?logo=go&logoColor=fff)](https://golang.org) [![Gin Framework](https://img.shields.io/badge/Gin-1.10.0-blue?logo=gin&logoColor=fff)](https://gin-gonic.com/) [![GORM](https://img.shields.io/badge/GORM-ORM-yellow)](https://gorm.io/) [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-blue?logo=postgresql&logoColor=fff)](https://www.postgresql.org/) [![Validator](https://img.shields.io/badge/Validator-v10.20.0-green)](https://github.com/go-playground/validator) [![Mapstructure](https://img.shields.io/badge/Mapstructure-v2.2.1-6f42c1)](https://github.com/go-viper/mapstructure) [![Wire](https://img.shields.io/badge/Google%20Wire-v0.6.0-blue?logo=google&logoColor=fff)](https://github.com/google/wire) [![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<br />
<br />
<div align="center">
  <a href="https://github.com/ramadiaz">
    <img src="https://go.dev/images/gophers/motorcycle.svg" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Xanny Go</h3>

  <p align="center">
    An awesome Golang API template to jumpstart your projects!
    <br />
    <a href="https://github.com/ramadiaz/xanny-go"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://xann.my.id/incognito">Report Bug</a>
    &middot;
    <a href="https://xann.my.id/incognito">Request Feature</a>
  </p>
</div>

## About The Project

**Xanny Go** is a robust and scalable boilerplate for building RESTful APIs using Go (Golang). Designed with best practices and modern development in mind, this template serves as a starting point for backend projects, enabling developers to quickly set up and customize their API solutions.

---

## Project Structure & Features

### Directory Structure

```
.
├── api
│   ├── blueprint
│   │   ├── controllers
│   │   │   ├── blueprint_ctrl.go
│   │   │   └── blueprint_ctrl_impl.go
│   │   ├── dto
│   │   │   └── response.go
│   │   ├── repositories
│   │   │   ├── blueprint_repo.go
│   │   │   └── blueprint_repo_impl.go
│   │   └── services
│   │       ├── blueprint_svc.go
│   │       └── blueprint_svc_impl.go
│   └── users
│       ├── controllers
│       │   ├── users_ctrl.go
│       │   └── users_ctrl_impl.go
│       ├── dto
│       │   ├── request.go
│       │   └── response.go
│       ├── repositories
│       │   ├── users_repo.go
│       │   └── users_repo_impl.go
│       └── services
│           ├── users_svc.go
│           └── users_svc_impl.go
├── cmd
│   ├── migrate
│   │   └── migrate.go
│   └── server
│       └── main.go
├── docs
│   ├── docs.go
│   ├── README.md
│   ├── swagger.json
│   └── swagger.yaml
├── emails
│   ├── dto
│   │   └── input.go
│   ├── services
│   │   └── emails_svc_impl.go
│   └── templates
│       └── example.html
├── injectors
│   ├── injector.go
│   └── wire_gen.go
├── internal
│   ├── auth
│   │   ├── controllers
│   │   │   ├── internal_auth_ctrl.go
│   │   │   └── internal_auth_ctrl_impl.go
│   │   ├── dto
│   │   │   ├── request.go
│   │   │   └── response.go
│   │   └── services
│   │       ├── internal_auth_svc.go
│   │       └── internal_auth_svc_impl.go
│   ├── injectors
│   │   ├── injector.go
│   │   └── wire_gen.go
│   └── routers
│       ├── auth_router.go
│       └── main_router.go
├── models
│   ├── clients_model.go
│   ├── refresh_token_model.go
│   └── users_model.go
├── pkg
│   ├── cache
│   │   ├── cache.go
│   │   ├── controller.go
│   │   ├── example.go
│   │   ├── factory.go
│   │   ├── integration.go
│   │   ├── memory_cache.go
│   │   ├── middleware.go
│   │   ├── README.md
│   │   ├── redis_cache.go
│   │   ├── router.go
│   │   └── service.go
│   ├── config
│   │   ├── database_config.go
│   │   └── env_config.go
│   ├── exceptions
│   │   ├── database_exception.go
│   │   ├── exception.go
│   │   └── variable.go
│   ├── helpers
│   │   ├── database_helper.go
│   │   ├── example_helper.go
│   │   ├── hash_helper.go
│   │   ├── health_helper.go
│   │   └── redis_helper.go
│   ├── logger
│   │   ├── general_log.go
│   │   └── startup_log.go
│   ├── mapper
│   │   └── users_mapper.go
│   ├── middleware
│   │   ├── auth_middleware.go
│   │   ├── cache_middleware.go
│   │   ├── gzip_middleware.go
│   │   ├── internal_middleware.go
│   │   ├── log_middleware.go
│   │   └── ratelimit_middleware.go
│   └── whatsapp
│       └── fonnte.go
├── routers
│   ├── main_router.go
│   └── users_router.go
├── LICENSE.txt
├── SWAGGER_SETUP.md
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── makefile
└── tmp/
```

### Main Features & Modules

#### 1. Modular API (Blueprint & Users)
- Each module (blueprint, users) consists of controller, service, repository, and DTO (Data Transfer Object).
- Example endpoints: user CRUD, login, refresh token, token blacklist, etc.

#### 2. Internal Auth
- Internal module for admin/internal authentication (internal/auth).
- Supports internal login, JWT validation, etc.

#### 3. Email Service
- Email sending with HTML template (emails/services, emails/templates).
- DTO for email input.

#### 4. Caching
- Supports memory cache & Redis (pkg/cache).
- Cache middleware, controller, service, and Redis integration.

#### 5. Middleware
- Authentication (auth_middleware.go)
- Rate Limiting (ratelimit_middleware.go)
- Logging (log_middleware.go)
- Gzip Compression (gzip_middleware.go)
- Internal Middleware, Cache Middleware

#### 6. Database & ORM
- User, client, refresh token models (models/)
- GORM ORM, auto-migration (cmd/migrate, pkg/config)
- Exception & helper for database

#### 7. Dependency Injection
- Using Google Wire (injectors/, internal/injectors/)

#### 8. API Documentation
- Swagger (docs/swagger.yaml, docs/swagger.json)
- Setup instructions in SWAGGER_SETUP.md

#### 9. Utilities & Helpers
- Hashing, health check, redis helper, etc (pkg/helpers)
- Logger (pkg/logger)
- Mapper (pkg/mapper)
- Exception handler (pkg/exceptions)

#### 10. Deployment & Automation
- Dockerfile & docker-compose.yaml
- Makefile for build, run, migrate, etc

#### 11. Routers
- Main router, users router, internal routers (routers/, internal/routers/)

#### 12. Others
- Modular structure, clean code, scalable, ready to use for REST API

---

### Key Features

- **Modular Architecture**: Organized with a clean directory structure, separating concerns for better maintainability and scalability.
- **Layered Approach**: Implements the handler-service-repository pattern to keep business logic, data access, and API handling separate.
- **Gin Framework**: Leverages the powerful and minimalist Gin framework for efficient request handling and routing.
- **ORM Support**: Uses GORM for database interactions, providing an easy and familiar way to handle models and migrations.
- **Configuration Management**: Centralized configuration management to simplify environment-specific settings using Viper.
- **JWT Authentication**: Secure authentication mechanism using JSON Web Tokens (JWT) for stateless and secure user sessions.
- **Custom Error Handling**: Centralized error management to ensure consistent and informative API responses.
- **Auto-Migrations**: Automated database migrations using a dedicated migration script for seamless schema updates.
- **Makefile Automation**: Includes a Makefile for common tasks such as running the application, building binaries, and executing migrations.

### Why Use Xanny Go?

This template is ideal for developers looking for a ready-to-use, yet flexible structure for their Go backend projects. By providing essential features out-of-the-box, it allows you to focus on building core functionalities without worrying about boilerplate code.


### Built With

Xanny Go is built using a combination of powerful tools and libraries that enable fast, scalable, and secure backend development. Here's a list of the technologies and frameworks used in this project:

* ![Go](https://img.shields.io/badge/GO-blue?style=for-the-badge&logo=go&logoColor=white)
* ![Gin Framework](https://img.shields.io/badge/Gin-blue?style=for-the-badge&logo=gin&logoColor=fff)
* ![GORM](https://img.shields.io/badge/GORM-yellow?style=for-the-badge)
* ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-blue?style=for-the-badge&logo=postgresql&logoColor=fff)
* ![Validator](https://img.shields.io/badge/Validator-green?style=for-the-badge)
* ![Mapstructure](https://img.shields.io/badge/Mapstructure-6f42c1?style=for-the-badge)
* ![Wire](https://img.shields.io/badge/Google%20Wire-blue?style=for-the-badge&logo=google&logoColor=fff)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Getting Started

To get a local copy of Xanny Go up and running on your machine, follow these simple steps.

### Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go (Golang)**: Version 1.24.2 or higher. You can download it from [Go's official website](https://golang.org/dl/).
- **Database**: You need a running instance of a relational database (e.g., PostgreSQL, MySQL, or MSSQL). Make sure to set up the database and configure the connection.
- **Make**: Optionally, you can install [Make](https://www.gnu.org/software/make/) for automating tasks like building, running migrations, etc.

### Installing

1. **Clone the repository**:
   Open your terminal and clone the project to your local machine using Git.

   ```bash
   git clone https://github.com/ramadiaz/xanny-go.git
   cd xanny-go
   ```

2. **Install dependencies**:
   Ensure your Go modules are set up and all dependencies are installed.

   ```bash
   go mod tidy
   ```

3. **Configure environment variables**:
   Set up your environment variables by creating a `.env` file in the root of the project. Here's an example:

   ```bash
   DB_USER=your-db-username
   DB_PASSWORD=your-db-password
   DB_HOST=your-db-host
   DB_PORT=your-db-port
   DB_NAME=your-db-name
   PORT=your-desire-port
   JWT_SECRET=your-jwt-secret
   
   ENVIRONMENT=production/development
   
   ADMIN_USERNAME=your-desire-username
   ADMIN_PASSWORD=your-desire-password
   
   REDIS_ADDR=your-redis-address
   REDIS_PASS=your-redis-password
   ```

   Replace the values with your own database connection details and secret key.

4. **Run database migrations**:
   The template comes with a migration script to set up the necessary database tables. You can run it using the following command:

   ```bash
   make migrate
   ```

5. **Run the application**:
   To start the application, you can use the following command:

   ```bash
   make run
   ```

   This will start the server, and the API will be accessible at `http://localhost:<PORT>`.

### Optional: Building the Application

If you'd like to build the application binary for production use, you can run:

```bash
make build
```

This will compile the application into a binary located in the `bin/` folder. You can then execute the binary directly.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---


## Usage

Once you have set up the project and run the application, you can access the default API health check endpoint to verify that the server is working correctly.

### Default API Endpoint

By default, the API exposes a health check endpoint:

- **GET** `/api/health`

This endpoint will return a JSON response with the health status of the application and its dependencies.

#### Example Request

```bash
curl http://localhost:<PORT>/api/health
```

Replace `<PORT>` with the port number specified in your `.env` file.

#### Example Response

```json
{
  "status": "healthy",
  "database": "healthy",
  "redis": "healthy"
}
```

This response indicates that the server and its dependencies are up and running. You can modify this route or add more routes as needed to expand the functionality of the API.

### Customizing the Server Port

To customize the port, change the `PORT` value in the `.env` file:

```bash
PORT=your-desired-port
```

Then, restart the application with:

```bash
make run
```

The API will now be accessible at `http://localhost:<your-desired-port>/api/health`.

---


## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

----


## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---


## Contact

Rama Diaz - [@ramadiazr](https://instagram.com/ramadiazr) - ramadiaz221@gmail.com

Project Link: [https://github.com/ramadiaz/xanny-go](https://github.com/ramadiaz/xanny-go)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
