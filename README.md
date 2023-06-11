# TSB API

## Description

The TSB API is a Go-based server application using the Fiber web framework. Its main function is to interact with the TSB API service, retrieve the URL of the most recently published Casco Excel file, download this file, and read its contents. The data from the Excel file, which includes brand, model, and Casco values, is then upserted (either created or updated) in a database.

This project allows users to save this data by simply changing the database settings in the environment (.env) file to match their own database configuration.

## Installation

1. Clone the repository
2. Run `cp .env.dist .env`
3. Update the .env file with your database settings
4. Run `go run server.go` or `go build server.go` and then `./server` to start the server
6. Send a POST request to the API endpoint (see below)
7. Check your database to see the results

### API Endpoint

The API has one route:

- POST `http://127.0.0.1:3000/api/get-files`: Accepts a JSON payload with a "year" field (integer only, e.g., 2023).

Example request:
```json
{
  "year": 2023
}
```
## Configuration

Database settings can be configured in the .env file.

## Models

The application works with the following data models:

- Brand
- Model
- CascoValue

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT license.
