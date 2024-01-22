## How to run

1. Clone the repository
1. Make sure to have Go
1. Install the dependencies running `go mod tidy`
1. Run the application with `go run .`
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
On your browser, go to ``127.0.0.1:8000/send`` to send a Push message and ``127.0.0.1:8000/receive`` to receive the messages. **Make sure to allow notifications!**

## Database
The application uses a PostgreSQL database to store data. The database connection is set up in the ``connection.go`` file. The database schema is defined in the ``models.go`` file.

Use ``scooters.sql`` and ``init.sql`` to load some test scooters and users.

## Firebase
The application uses Firebase for real-time updates. The Firebase configuration is set up in the ``firebase.go`` file. The Firebase Cloud Messaging service worker is defined in the ``firebase-messaging-sw.js`` file.