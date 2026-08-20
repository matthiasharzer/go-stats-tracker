# Pokémon Go Stats Tracker

A personal user XP tracker based on screenshots using Tesseract.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<br>

## Setup

### Create a Google Cloud OAuth 2.0 Client

This project connects to the Google Sheets API using OAuth 2.0. To use it, you need to create a Google Cloud project and set up an OAuth 2.0 client.
1. Create a project in the [Google Cloud Console](https://console.cloud.google.com/).
2. Navigate to "APIs & Services" > "Credentials" and click "Create Credentials" > "OAuth client ID".
3. Select "Web application" as the application type and set the redirect URI to `<YOUR_REDIRECT_URL>/api/v1/callback`.
4. Note down the `Client ID` and `Client Secret` for later use.
5. Enable the Google Sheets API for your project by navigating to "APIs & Services" > "Library" and searching for "Google Sheets API". Click "Enable".
6. Add the `https://www.googleapis.com/auth/spreadsheets` scope to your OAuth consent screen by navigating to "APIs & Services" > "OAuth consent screen" > "Data Access".

### Docker (recommended)

The easiest way to run the application is with Docker. A pre-built image is available on the [GitHub Container Registry](https://ghcr.io/matthiasharzer/go-stats-tracker).


Create an `.env` file with the following content, using the `Client ID`, `Client Secret`, and redirect URL you obtained from the Google Cloud Console:
```env
REDIRECT_URL=<YOUR_REDIRECT_URL>/api/v1/callback
CLIENT_ID=<YOUR_CLIENT_ID>
CLIENT_SECRET=<YOUR_CLIENT_SECRET>
```

Run the following command to start the application on port 4000:
```bash
docker run --env-file .env -p 4000:4000 ghcr.io/matthiasharzer/go-stats-tracker:latest run -p 4000 -d data/db.sqlite3
```

> [!IMPORTANT]
> When mounting a directory, make sure it has the correct permission for the docker container user: `chown -R 1000:1000 <your-host-mount-path>`. This is due to the app not running as root inside the container.


#### Docker Compose
You can use the [`docker-compose.yml`](./docker-compose.yml) file to start the application with Docker Compose. Run the following command to start the application on port 4000:
```bash
docker compose up -d
```
> Make sure a `.env` file is present in the same directory as the `docker-compose.yml` file and the `./data` directory exists and has the correct permissions for the docker container user: `chown -R 1000:1000 ./data`

### Binary

You can download the [latest release](https://github.com/matthiasharzer/go-stats-tracker/releases/latest) from GitHub. Prebuilt binaries are available for Linux (x86_64). Other platforms may be supported in the future.


## Usage

### Check The Version
Check the version of the application:
```bash
go-stats-tracker version
```

### Run The Server
Start the server using:
```bash
go-stats-tracker run [-p <port>] [--host <host>] [-d <db-file-path>]
```

Requires the following environment variables to be set:
- `REDIRECT_URL`: The redirect URL for the OAuth 2.0 client. This must match exactly the URL configured in the Google Cloud Console.
- `CLIENT_ID`: The client ID for the OAuth 2.0 client.
- `CLIENT_SECRET`: The client secret for the OAuth 2.0 client.

Command line arguments:

| Argument                 | Required | Default                 | Description                                                      |
|--------------------------|----------|-------------------------|------------------------------------------------------------------|
| `--port` / `-p`          | ❌       | `4000`                  | The port to start the HTTP-server on                             |
| `--host`                 | ❌       | `""` _(all interfaces)_ | The host to start the HTTP-server on                             |
| `--database-file` / `-d` | ❌       | `data/db.sqlite`        | The location of the database file where user context is saved in |	

### Server Endpoints

The server exposes the following endpoints:

| Endpoint                   | Method | Params / Payload    | Description                                                                                                                                                             |
|----------------------------|--------|---------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `/api/v1/health`           | GET    |                     | Check the health of the server                                                                                                                                          |
| `/api/v1/register`         | GET    | `sheet_id`          | Register a new sheet to populate. Will trigger the OAuth flow to authenticate with your Google account. Returns a unique user ID which is used to authenticate later on |
| `/api/v1/callback`         | GET    |                     | The OAuth callback endpoint                                                                                                                                             |
| `/api/v1/ingest/{user-id}` | POST   | _(screenshot data)_ | Submits a new screenshot. Authentication is handled by the user ID which resolves to the prior configured sheet ID                                                      |


### Spreadsheet Template

The application currently expects a specific spreadsheet template to be used. You can find the template [here](https://docs.google.com/spreadsheets/d/1ls99T0X2Dy3nMDzR695GSkmfhdjUaF--T71CFu6zNQg/edit).

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details
