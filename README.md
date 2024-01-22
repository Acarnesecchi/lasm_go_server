## How to run

1. Clone the repository
1. Build or pull the container image
1. Run a container with the image and define a volume `docker run -n api-server -v crentials:app/credentials api-server` or `docker run -n login-server -v crentials:app/credentials login-server`
> It will execute a RESTful API on port 8080 and a web service on port 8000.

## Configuration
This project requires two configuration files to run:

**postgres.ini:** This file contains the configuration for the PostgreSQL database connection. It should be located in the root directory of the project and have the following structure:
```ini
[postgresql]
host = localhost
port = 5432
user = yourusername
password = yourpassword
database = yourdatabase
```
**serviceAccountKey.json:** This file contains the Firebase service account key, which is used to authenticate the server with Firebase. It should be located in the credentials directory of the project.


## API Endpoints
The application provides several API endpoints for managing users, scooters, and rentals. These endpoints are defined in the ``endpoints.go`` file.

## Web Services
On your browser, go to ``127.0.0.1:8000/html`` to see the firebase SSO login page.

## Database
The application uses a PostgreSQL database to store data. The database connection is set up in the ``connection.go`` file. The database schema is defined in the ``models.go`` file.

Use ``scooters.sql`` and ``init.sql`` to load some test scooters and users.

## Testing
To test the different endpoints either use Postman or curl:

· `http://localhost:8080/endpoints/login`: Send the login information in the body as `raw`.
```json
{
    "Username": "alex",
    "Password": "alex1234"
}
```

· `http://localhost:8080/endpoints/scooter/<UUID>` and `http://localhost:8080/endpoints/rent/ACTION/<UUID>`: Include either the scooter or rent UUID in the URL