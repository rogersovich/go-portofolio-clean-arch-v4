# Portfolio Backend Application

[](https://golang.org/)

This backend application serves as the API foundation for my personal portfolio, providing structured and managed data for various sections of the portfolio's frontend. Built with **Go** and utilizing the **Gin** framework, this application is designed for performance, scalability, and ease of maintenance.

-----

## Table of Contents

- [Key Features](#key-features)
- [Technologies Used](#technologies-used)
- [Modules](#modules)
- [System Requirements](#system-requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Contributing](#contributing)
- [License](#license)

-----

## Key Features

  * **RESTful API:** Provides clean and structured API endpoints for frontend data consumption.
  * **Secure Authentication:** Uses JWT (JSON Web Tokens) for secure user authentication.
  * **Data Validation:** Ensures data integrity with comprehensive input validation.
  * **Object Storage:** Integration with MinIO for file and media management.
  * **Centralized Logging:** Utilizes Logrus for effective application logging.
  * **Database Management:** GORM ORM for seamless interaction with MySQL.

-----

## Technologies Used

  * **Go (Golang):** The core programming language for performance and concurrency.
  * **Gin:** A high-performance web framework for building RESTful APIs.
  * **GORM:** An ORM (Object-Relational Mapping) for Go, simplifying MySQL database interactions.
  * **JWT Go:** JSON Web Tokens implementation for authentication.
  * **MinIO:** High-performance S3 compatible object storage server.
  * **MySQL:** Relational database management system.
  * **go-playground/validator:** A flexible Go validation library.
  * **godotenv:** Loads environment variables from a `.env` file.
  * **Logrus:** A powerful and flexible logging library for Go.

-----

## Modules

This backend application consists of the following modules, each managing specific functionality:

  * **`about`:** Manages "About Me" information for the portfolio.
  * **`auth`:** Handles user authentication and authorization processes.
  * **`author`:** Manages author details (if multiple authors for the blog).
  * **`blog`:** Manages blog posts, including CRUD (Create, Read, Update, Delete) and related functionalities.
  * **`experience`:** Stores and manages work or education experience details.
  * **`project`:** Manages information about completed projects.
  * **`reading_time`:** Calculates and stores estimated reading time for blog posts.
  * **`statistic`:** Collects and manages statistics related to portfolio usage (e.g., visit count).
  * **`technology`:** Manages a list of technologies used (e.g., Vue, React, Go).
  * **`testimonial`:** Manages testimonials or reviews.
  * **`topic`:** Manages topics or categories for blog posts.
  * **`user`:** Manages user information, including profiles and roles.

-----

## System Requirements

Ensure you have the following installed on your system:

  * **Go** (Version 1.18 or higher)
  * **MySQL** (Version 5.7 or higher)
  * **MinIO** (Server or access to a MinIO service)
  * **Air** (for live reloading during development)

-----

## Installation

Follow these steps to install and run the application locally:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/rogersovich/go-portofolio-clean-arch-v4
    cd your-repo-name
    ```

2.  **Install Go dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Install Air (if you haven't already):**

    ```bash
    go install github.com/cosmtrek/air@latest
    ```

4.  **Set up MySQL database:**

      * Create a new database for this application (e.g., `portfolio_db`).
      * Ensure the database user has appropriate access rights.

5.  **Set up MinIO:**

      * Ensure your MinIO server is running and you have access credentials (access key and secret key).
      * Create a bucket to be used for file storage.

-----

## Configuration

Create a `.env` file in your project root based on `.env.example` and fill in your configuration details:

```dotenv
# .env example
APP_PORT=YOUR_PORT
APP_ENV=production

# MySQL on the host
DB_HOST=
DB_PORT=
DB_NAME=
DB_USER=
DB_PASSWORD=

# MinIO in Docker
MINIO_ENDPOINT_UPLOAD=
MINIO_ENDPOINT_VIEW=
MINIO_KEY_ID=
MINIO_KEY_SECRET=
MINIO_SSL=
MINIO_BUCKET=

JWT_SECRET=
```

**Note:** Ensure your `JWT_SECRET_KEY` is a strong and unique string.

-----

## Running the Application

After configuration, you can run the application using **Air** for live reloading during development:

1.  **Migrate Database (optional, if you have GORM migrations):**
    This application uses GORM, which can automatically manage database schema migrations. If you prefer to run migrations manually or have specific migrations, you might need to add commands here (e.g., `go run main.go migrate`).

2.  **Run with Air:**

    ```bash
    air
    ```

    Air will automatically restart the application whenever you make changes to your Go files. The application will be running at `http://localhost:APP_PORT` (default port 4000).

-----

## Contributing

Contributions are highly appreciated\! If you have ideas for improvements, new features, or want to report a bug, please create an *issue* or submit a *pull request*.

-----

## License

This project is licensed under the [MIT License](LICENSE).

-----