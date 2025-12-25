# Ergometer.Live REST API Server

REST API server for the Ergometer.Live application providing authentication, workout management, and data persistence.

## Features

- Firebase Authentication with Google OAuth
- InfluxDB integration for workout time-series data
- CORS support for frontend communication
- Request logging middleware
- Health check endpoint

## Prerequisites

- Go 1.23 or later
- Firebase project with service account credentials
- InfluxDB Cloud account (or self-hosted InfluxDB)

## Setup

### 1. Install Dependencies

```bash
cd api
go mod download
```

### 2. Configure Environment Variables

Copy the example environment file and update with your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your Firebase and InfluxDB credentials:

```bash
PORT=3000
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_CREDENTIALS_PATH=/path/to/service-account.json
INFLUXDB_URL=https://your-influxdb-url.com
INFLUXDB_TOKEN=your-token
INFLUXDB_ORG=your-org
INFLUXDB_BUCKET=ergometer-workouts
ALLOWED_ORIGINS=http://localhost:5173
```

### 3. Firebase Setup

1. Create a Firebase project at https://console.firebase.google.com
2. Enable Google Authentication in Firebase Console
3. Download service account JSON:
   - Go to Project Settings > Service Accounts
   - Click "Generate New Private Key"
   - Save the JSON file securely
   - Update `FIREBASE_CREDENTIALS_PATH` in `.env`

### 4. InfluxDB Setup

1. Create an InfluxDB Cloud account at https://cloud2.influxdata.com
2. Create a new bucket named `ergometer-workouts`
3. Generate an API token with read/write access
4. Update InfluxDB configuration in `.env`

## Running the Server

### Development Mode

```bash
go run main.go
```

The server will start on `http://localhost:3000` (or the port specified in `.env`).

### Production Build

```bash
go build -o ergometer-api
./ergometer-api
```

## API Endpoints

### Health Check

```
GET /health
```

Returns server health status. No authentication required.

**Response:**
```json
{
  "status": "healthy"
}
```

### Authentication

```
POST /api/v1/auth/verify
```

Verifies Firebase ID token. Requires Authorization header.

**Headers:**
```
Authorization: Bearer <firebase-id-token>
```

**Response:**
```json
{
  "uid": "user-id",
  "verified": true
}
```

## Middleware

### Authentication Middleware

All endpoints under `/api/v1/auth/*` require a valid Firebase ID token in the Authorization header:

```
Authorization: Bearer <firebase-id-token>
```

### CORS Middleware

Configured to allow requests from origins specified in `ALLOWED_ORIGINS` environment variable.

### Logger Middleware

Logs all HTTP requests with method, path, status code, and duration.

## Project Structure

```
api/
├── main.go              # Server entry point
├── config/
│   └── config.go       # Environment configuration
├── middleware/
│   ├── auth.go         # Firebase token validation
│   ├── cors.go         # CORS configuration
│   └── logger.go       # Request logging
├── services/
│   ├── firebase.go     # Firebase Admin SDK
│   └── influxdb.go     # InfluxDB client (planned)
├── handlers/
│   ├── auth.go         # Auth endpoints (planned)
│   ├── workouts.go     # Workout CRUD (planned)
│   └── views.go        # Widget layouts (planned)
└── models/
    ├── user.go         # User model (planned)
    ├── workout.go      # Workout model (planned)
    └── view.go         # View model (planned)
```

## Testing

Test the health check endpoint:

```bash
curl http://localhost:3000/health
```

Test authentication (requires valid Firebase token):

```bash
curl -X POST http://localhost:3000/api/v1/auth/verify \
  -H "Authorization: Bearer <your-firebase-token>"
```

## Next Steps

- Implement InfluxDB service for workout data storage
- Add workout CRUD endpoints
- Add widget layout (views) endpoints
- Implement PostgreSQL for user metadata (optional)
- Add rate limiting
- Add request validation
- Write unit tests

## Security Notes

- Never commit `.env` file or Firebase credentials to version control
- Restrict CORS origins in production
- Use HTTPS in production
- Implement rate limiting for production
- Validate all user inputs
- Keep dependencies updated
