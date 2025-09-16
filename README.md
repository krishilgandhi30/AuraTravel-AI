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
```

4. **Run Backend**
```bash
go run .
# Or build and run
go build -o server .
./server
```

The backend will be available at `http://localhost:8080`

## 📊 API Documentation

### Health Check
- `GET /health` - Server health status

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh
- `POST /api/v1/auth/firebase-auth` - Firebase authentication

### Trip Management
- `POST /api/v1/trips/` - Create new trip
- `GET /api/v1/trips/` - Get user trips
- `GET /api/v1/trips/:id` - Get specific trip
- `PUT /api/v1/trips/:id` - Update trip
- `DELETE /api/v1/trips/:id` - Delete trip
- `GET /api/v1/trips/recommendations` - Get AI recommendations

### AI-Powered Features
- `POST /api/v1/ai/plan-trip` - Generate AI trip plan
- `GET /api/v1/ai/recommendations` - Get AI recommendations
- `POST /api/v1/ai/optimize/:id` - Optimize existing itinerary
- `POST /api/v1/ai/analyze-image` - Analyze destination images
- `GET /api/v1/ai/insights` - Get travel insights

### User Management
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile

## 🤖 AI Features

### Gemini AI Integration
- Intelligent trip planning and itinerary generation
- Natural language processing for user preferences
- Dynamic recommendations based on user input

### Vertex AI Analytics
- Advanced destination analysis
- Travel pattern predictions
- Personalized recommendations engine

### Cloud Vision API
- Destination image analysis
- Visual content recommendations
- Automated tagging and categorization

### BigQuery Analytics
- Travel data analysis and insights
- Budget optimization recommendations
- Historical travel pattern analysis

## 🔧 Development

### Running Tests
```bash
# Frontend tests
cd frontend
npm test

# Backend tests
cd backend
go test ./...
```

### Building for Production
```bash
# Frontend
cd frontend
npm run build

# Backend
cd backend
go build -o server .
```

## 🌐 Environment Variables

### Backend (.env)
```env
PORT=8080
ENVIRONMENT=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=auratravel_db
JWT_SECRET=your_jwt_secret
GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account.json
GEMINI_API_KEY=your_gemini_api_key
FIREBASE_PROJECT_ID=your_firebase_project_id
# ... (see .env file for complete list)
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_FIREBASE_API_KEY=your_firebase_api_key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_project_id
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 📈 Features

### Current Features
- ✅ User authentication (Firebase + JWT)
- ✅ Trip creation and management
- ✅ AI-powered trip planning with Gemini
- ✅ Destination recommendations
- ✅ Image analysis for destinations
- ✅ Real-time travel insights
- ✅ Budget analysis and optimization
- ✅ Responsive web interface

### Upcoming Features
- 🔄 Real-time collaborative trip planning
- 🔄 Integration with booking platforms
- 🔄 Mobile app (React Native)
- 🔄 Advanced analytics dashboard
- 🔄 Social features and trip sharing

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Google Cloud AI Platform for powerful AI capabilities
- Firebase for authentication and real-time features
- Next.js and Go communities for excellent frameworks
- Open source contributors who made this project possible