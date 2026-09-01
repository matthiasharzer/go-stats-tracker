# Pokémon Go Stats Tracker

A personal user XP tracker based on screenshots using Tesseract.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<br>

## Setup

### Create a Google Cloud OAuth 2.0 Client

This project connects to the Google Sheets API using OAuth 2.0. To use it, you need to create a Google Cloud project and set up an OAuth 2.0 client.
1. Create a project in the [Google Cloud Console](https://console.cloud.google.com/).
2. Navigate to "APIs & Services" > "Credentials" and click "Create Credentials" > "OAuth client ID".
3. Select "Web application" as the application type and set the redirect URI to `<YOUR_REDIRECT_URL>/api/v1/auth/callback`.
4. Note down the `Client ID` and `Client Secret` for later use.
5. Enable the following APIs, navigating to "APIs & Services" > "Library":
   - [Google Sheets API](https://console.cloud.google.com/marketplace/product/google/sheets.googleapis.com)
   - [Google Drive API](https://console.cloud.google.com/marketplace/product/google/drive.googleapis.com)
   - [Google Picker API](https://console.cloud.google.com/marketplace/product/google/picker.googleapis.com)
6. Add the `https://www.googleapis.com/auth/drive.file` scope to your OAuth consent screen by navigating to "APIs & Services" > "OAuth consent screen" > "Data Access".

### Create an API key for the Google Picker API 

1. Navigate to "APIs & Services" > "Credentials" and click "Create Credentials" > "API Key".
2. Select the "Google Picker API".
3. (Optionally) Further restrict the API key.
4. Click create and note down the API key for later

### Gather the App ID

 The app ID is required for the Google Picker API and is part of the authorization chain to access a user's files

1. Navigate to "IAM & Admin" > "Settings"
2. Copy the "Project number", which is referred to as the App ID


### Docker (recommended)

The easiest way to run the application is with Docker. A pre-built image is available on the [GitHub Container Registry](https://ghcr.io/matthiasharzer/go-stats-tracker).


Create an `.env` file with the following content, using the `Client ID`, `Client Secret`, `Picker API Key`, `App ID`, and redirect URL you obtained from the Google Cloud Console:
```env
REDIRECT_URL=<YOUR_REDIRECT_URL>/api/v1/auth/callback
CLIENT_ID=<YOUR_CLIENT_ID>
CLIENT_SECRET=<YOUR_CLIENT_SECRET>
APP_ID=<YOUR_APP_ID>
PICKER_API_KEY=<YOUR_FILE_PICKER_API_KEY>
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
- `APP_ID`: The unique [app ID from the GCP-project](https://console.cloud.google.com/iam-admin/settings), where it is referred to as "Project number"
- `PICKER_API_KEY`: The API key used by the Google Picker API.

Command line arguments:

| Argument                 | Required | Default                 | Description                                                      |
|--------------------------|----------|-------------------------|------------------------------------------------------------------|
| `--port` / `-p`          | ❌       | `4000`                  | The port to start the HTTP-server on                             |
| `--host`                 | ❌       | `""` _(all interfaces)_ | The host to start the HTTP-server on                             |
| `--database-file` / `-d` | ❌       | `data/db.sqlite`        | The location of the database file where user context is saved in |	

### Server Endpoints

The server exposes the following endpoints:

| Endpoint                    | Method | Params / Payload    | Authentication       | Description                                                                                                                                                                 |
|-----------------------------|--------|---------------------|----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `/api/v1/health`            | GET    |                     |                      | Check the health of the server                                                                                                                                              |
| `/api/v1/submit/{sheet_id}` | POST   | _(screenshot data)_ | Bearer authenication | Submits a new screenshot and populates the spreadsheet (`sheet_id`). Authentication is handled by the access token provided via the `Authorization: Bearer <token>` header. |
| `/api/v1/auth/register`     | GET    |                     |                      | Register a new sheet to populate. Will trigger the OAuth flow to authenticate with your Google account. Returns a unique user ID which is used to authenticate later on     |
| `/api/v1/auth/callback`     | GET    |                     |                      | The OAuth callback endpoint. Will be redirected to after authenticating with Google when registering.                                                                       |
| `/api/v1/auth/link`         | POST   |                     |                      | Links a user selected spreadsheet.                                                                                                                                          |

To use the tool, go the `/api/v1/auth/register` URL in your browser and follow the steps to authorize access to a spreadsheet.

### Spreadsheet Template

The application currently expects a specific spreadsheet template to be used. You can find the template [here](https://docs.google.com/spreadsheets/d/1ls99T0X2Dy3nMDzR695GSkmfhdjUaF--T71CFu6zNQg/edit).

### Tracking Screenshots (Android only)

To automatically submit screenshots made from the Pokémon Go app, you can use [Automate](https://llamalab.com/automate/) to create a flow that detects new screenshots and sends them to the server. 
The [AutomateTemplate.flo](./AutomateTemplate.flo) file contains a sample flow that may be used as a starting point. You will need to adjust the flow depending on your device/android capabilities and need to configure the server URL and user ID in the flow.

There may be other ways to automatically submit screenshots, but this is the only method that has been tested and confirmed to work.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details
