# AuraTravel AI - Full Stack Travel Planning Application

A comprehensive AI-powered travel planning application built with Next.js frontend and Go backend, leveraging Google AI technologies for intelligent trip recommendations and planning.

## 🏗️ Project Structure

```
AuraTravel-AI/
├── frontend/                 # Next.js 14 Frontend
│   ├── src/
│   │   ├── app/             # App Router pages
│   │   ├── components/      # Reusable UI components
│   │   └── lib/            # Utility functions and configs
│   ├── public/             # Static assets
│   └── package.json
│
├── backend/                 # Go Backend with Google AI
│   ├── cmd/                # Main application entry
│   ├── internal/           # Private application code
│   │   ├── config/        # Configuration management
│   │   ├── database/      # Database connection and setup
│   │   ├── handlers/      # HTTP request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   ├── models/        # Database models
│   │   ├── routes/        # API routing
│   │   └── services/      # Business logic and AI services
│   ├── .env               # Environment variables
│   ├── go.mod             # Go dependencies
│   └── main.go            # Application entry point
```

## 🚀 Tech Stack

### Frontend
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: Custom components with Lucide React icons
- **Animation**: Framer Motion
- **Authentication**: Firebase Auth integration
- **State Management**: React hooks and context

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: Firebase Admin SDK + JWT
- **AI Services**: Google AI Platform integration
  - **Gemini AI**: Trip planning and recommendations
  - **Vertex AI**: Advanced analytics and insights
  - **Cloud Vision**: Image analysis for destinations
  - **BigQuery**: Travel data analytics
  - **Firebase**: User management and real-time data

## 🛠️ Installation & Setup

### Prerequisites
- Node.js 18+ and npm/yarn
- Go 1.21+
- PostgreSQL 12+
- Google Cloud Project with AI APIs enabled

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:3000`

### Backend Setup

1. **Install Dependencies**
```bash
cd backend
go mod download
```

2. **Environment Configuration**
Copy `.env` file and update with your credentials:
```bash
cp .env .env.local
# Edit .env.local with your actual API keys and database credentials
```

3. **Database Setup**
```bash
# Install PostgreSQL and create database
createdb auratravel_db

# The application will auto-migrate tables on startup
# AuraTravel AI

AuraTravel AI is an end-to-end, AI-powered travel planning platform. It combines a Next.js frontend with a Go (Gin) backend and leverages Google Cloud AI (Gemini, Vertex AI), Cloud Vision, Firestore, and Firebase for authentication and data storage. The system uses a RAG (Retrieval Augmented Generation) approach with a vector database for personalized, real-time itinerary generation.

This README documents how to set up the project locally, architecture notes, where to find core components, environment variables, how to run and test, API endpoints, and special notes about Gemini/RAG/vector DB behavior.

---

## Table of Contents
- [Project layout](#project-layout)
- [Architecture overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Local setup (frontend & backend)](#local-setup-frontend--backend)
  - [Frontend](#frontend)
  - [Backend](#backend)
- [Environment variables](#environment-variables)
- [API reference (selected endpoints)](#api-reference-selected-endpoints)
- [Gemini, RAG, and Vector DB notes](#gemini-rag-and-vector-db-notes)
- [Testing](#testing)
- [Troubleshooting & common issues](#troubleshooting--common-issues)
- [Development tips & next steps](#development-tips--next-steps)
- [Contributing](#contributing)
- [License](#license)

---

## Project layout

Top-level layout (relevant folders):

```
AuraTravel-AI/
├── frontend/                 # Next.js 14 frontend (TypeScript + Tailwind)
├── backend/                  # Go backend (Gin) and AI services
│   ├── internal/
│   │   ├── config/           # configuration loader
│   │   ├── handlers/         # HTTP handlers (ai, trips, vector, etc.)
│   │   ├── routes/           # route registration
│   │   └── services/         # Gemini, RAG retriever, Vector DB, Vertex, etc.
│   └── main.go               # backend entrypoint
├── README.md                 # (this file)
```

## Architecture overview

- Frontend: Next.js (App Router) TypeScript application that calls backend REST APIs.
- Backend: Gin-based HTTP API written in Go. The backend orchestrates RAG retrieval, vector similarity search, embedding generation, and NLG via Gemini.
- Vector DB: Firestore-backed document store used to persist embeddings and metadata. The code contains a `VectorDatabase` service that stores embeddings and computes cosine similarity searches.
- Embeddings: Primary production path uses Vertex AI text embedding models (`textembedding-gecko`) when available; fallback to deterministic mock embeddings if credentials are missing.
- Gemini: Used for natural language generation (itinerary creation, recommendations). If `GEMINI_API_KEY` is not set, the backend falls back to mock generator implementations.

## Prerequisites

- Node.js 18+ and npm (or yarn)
- Go 1.21+
- PostgreSQL (or your preferred DB) for main app data (the repo contains DB wiring, check `internal/database`)
- Google Cloud project with:
  - Vertex AI API
  - Generative Language API (Gemini)
  - Firestore (if using vector storage)
  - Service account with required permissions

If you don't have Google Cloud credentials, the system provides mock fallbacks for Gemini and embeddings so you can run and test locally.

---

## Local setup (frontend & backend)

All terminal commands below assume PowerShell on Windows. Use your terminal of choice and adapt commands for bash if needed.

### Frontend

1. Install dependencies and run dev:

```powershell
cd frontend
npm install
npm run dev
```

2. Open `http://localhost:3000` in your browser.

Available scripts (in `frontend/package.json`):

- `dev` – run Next.js dev server
- `build` – build for production
- `start` – run built app
- `lint` – run linting

### Backend

1. Install Go dependencies

```powershell
cd backend
go mod download
```

2. Environment: copy `.env.example` or `.env` (if present) and set values. Example steps (PowerShell):

```powershell
cd backend
copy .env .env.local
# Then edit .env.local in your editor and fill values
```

3. Run the backend (development):

```powershell
cd backend
go run .
```

4. The backend listens on `:8080` by default. Health endpoint:

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
```

5. To build a production binary:

```powershell
cd backend
go build -o server .
.\server.exe
```

---

## Environment variables

Below are the main variables used by backend and frontend. See `internal/config/config.go` for full loader and additional keys.

### Backend (.env)

```
PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=auratravel_db

# Auth
JWT_SECRET=your_jwt_secret

# Google cloud
GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account.json
GEMINI_API_KEY=your_gemini_api_key

# Firebase (if used)
FIREBASE_PROJECT_ID=your_firebase_project_id
```

Notes:
- If `GEMINI_API_KEY` is empty, the backend will use mock Gemini implementations for development.
- If Google Cloud credentials are missing, embedding generation falls back to deterministic mock embeddings.

### Frontend (.env.local)

```
NEXT_PUBLIC_FIREBASE_API_KEY=your_firebase_api_key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_project_id
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## API reference (selected endpoints)

Prefix: `/api/v1`

- `GET /health` – health check (top-level, not under /api/v1 in `main.go`)

Authentication-protected endpoints (use Bearer token):

Trip & AI endpoints

- `POST /api/v1/trips/` – create a new trip
- `GET /api/v1/trips/` – list trips
- `GET /api/v1/trips/:id` – get trip
- `POST /api/v1/ai/plan-trip` – RAG + Gemini itinerary generation (main AI endpoint)
- `GET /api/v1/ai/recommendations` – destination recommendations

Vector / RAG endpoints

- `GET /api/v1/ai/rag-context` – get real-time RAG context for a destination (query params)
- `POST /api/v1/ai/validate-availability` – validate availability for items
- `POST /api/v1/vector/search-attractions` – search similar attractions using vector DB
- `POST /api/v1/vector/search-trips` – search similar trip embeddings
- `POST /api/v1/vector/store-preferences` – store user preferences as embedding
- `GET /api/v1/vector/predict-cost` – lightweight cost prediction API

For full handler behavior, see `backend/internal/handlers`.

---

## Gemini, RAG, and Vector DB notes

- Gemini (Generative Language API): the project includes a `GeminiService` that will call the Generative Language API when `GEMINI_API_KEY` is set. Otherwise the service returns deterministic mock responses useful for local development.
- RAG Retriever: aggregates context from multiple data sources (attractions, hotels, weather, events, user history). The retriever is implemented in `backend/internal/services/rag_retriever.go`.
- Vector DB: a Firestore-backed document collection stores embeddings and metadata. The `VectorDatabase` service generates embeddings using the `EmbeddingService` (Vertex AI) when configured. If embeddings cannot be generated, code falls back to a deterministic mock embedding generator to allow offline testing.

Security note: Do not commit production API keys or service account JSON files to the repository. Use environment variables or a secrets manager.

---

## Testing

- Backend unit tests: from the `backend` folder run:

```powershell
cd backend
go test ./...
```

- Frontend tests: the project currently focuses on building UI and integration; any test commands are in `frontend/package.json` and can be run with `npm test` if configured.

---

## Troubleshooting & common issues

- Health check fails:
  - Ensure backend is running on the expected port and `.env` variables are loaded.
  - Use `Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET` in PowerShell.

- Gemini calls failing:
  - If `GEMINI_API_KEY` is missing or invalid, the app uses mock Gemini functions. To enable real Gemini calls, set `GEMINI_API_KEY` and ensure your Google IAM permissions are correct.

- Embeddings not generated:
  - Ensure `GOOGLE_APPLICATION_CREDENTIALS` points to a valid service account JSON with Vertex AI access. Without it, the code uses a mock embedding generator.

- Firestore operations failing:
  - Ensure `FIREBASE_PROJECT_ID` and service account credentials are set and Firestore is initialized in your GCP project.

---

## Development tips & next steps

- To test the RAG flow end-to-end, run backend & frontend locally, toggle RAG mode in the planner UI, and inspect network requests to `/api/v1/ai/plan-trip` and `/api/v1/ai/rag-context`.
- Add integration tests for the `GeminiService` and `VectorDatabase` that mock external API calls.
- For production deploy, containerize the backend and use a CI pipeline to run tests, build the binary, and deploy to your cloud provider.

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please follow the repo's linting and testing rules and include unit tests for new backend functionality.

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

If you want, I can also:

- add an example `.env.example` file at the root
- generate a quick Postman collection for the main endpoints
- add a short developer checklist for running the RAG + Gemini end-to-end

Tell me which of these you'd like next and I will add them.