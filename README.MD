# LensPix

LensPix is a photo gallery web application built with Go. It allows users to sign up and create their photo galleries. The project focuses on user authentication, csrf protection, and seamless deployment using Docker and Make commands. This site is hosted on AWS but users can also deploy it locally by following the below instructions.

 **Live Link** : https://lenspixx.ayushsharma.co.in


## Features

- User Authentication: Sign up, log in, and create photo galleries.
- CSRF Protection: Enhances security with cross-site request forgery protection.
- Password Reset: Utilizes Mailtrap for password reset functionality.
- Database: Stores data in PostgreSQL for robust data management.
- Webpage Design: Uses Go HTML library and Tailwind CSS for a responsive UI.
- Session Tokens: Simplifies user sign-in with session tokens.
- Easy Deployment: Deploy LensPix effortlessly using Docker, with a multistage Docker build for efficiency. Hosted live on AWS ,users can setup locally also.
- Schema Migrations: Utilizes Goose for database schema migrations.

## Libraries Used

- **Chi:** A lightweight, idiomatic, and composable router for building Go HTTP services.
- **gorilla/csrf:** Provides Cross-Site Request Forgery (CSRF) protection middleware.
- **jackc/pgx:** A PostgreSQL driver for Go's database/sql package.
- **pressly/goose:** A database migration tool, handles versioning and migrations.
- **x/crypto:** A collection of packages for secure cryptography operations.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your machine.
- [Docker](https://www.docker.com/) for containerized deployment.

## Deploying the project locally

1. Clone the repository: `git clone https://github.com/ayushthe1/lenspix.git`
2. Navigate to the project directory: `cd lenspix`
3. Users need to provide a `.env` file for configuring the application. An `env.template` file has   been provided as a reference
3. Build and run the application: `make run`
4. Open your browser and visit [http://localhost:3000](http://localhost:3000)


## Contributing

All contributions are welcome. Just fork this repo and raise any issue or pull request .
