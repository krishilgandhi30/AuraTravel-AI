# 🌟 AuraTravel AI - Simple Version

A lightweight, standalone version of the Personal AI Travel Assistant that works without external dependencies.

## ✨ Features

- **🤖 AI Travel Assistant**: Chat with AI for personalized travel recommendations
- **✈️ Trip Planning**: Plan trips to amazing Indian destinations
- **🏔️ Multiple Destinations**: Rishikesh, Goa, Shimla, Manali and more
- **🎯 Vibe-based Matching**: Select your preferred travel vibes
- **💰 Budget Planning**: Set your budget and get appropriate suggestions
- **📱 Responsive Design**: Works on desktop and mobile

## 🚀 Quick Start

### Prerequisites
- Go (any recent version)
- A web browser

### Running the Application

1. **Start the Backend**:
   ```powershell
   # Windows
   .\start-simple.ps1
   
   # OR manually
   go run simple-backend.go
   ```

2. **Open the Frontend**:
   - Open `simple-frontend.html` in your web browser
   - The frontend will automatically connect to the backend at `localhost:8080`

3. **Start Using**:
   - Chat with the AI assistant for travel recommendations
   - Use the trip planner to create personalized itineraries
   - Explore different destinations and vibes

## 📋 API Endpoints

The backend provides the following REST API endpoints:

- `GET /health` - Health check
- `GET /api/v1/destinations` - Get all destinations
- `GET /api/v1/destinations/search?q=query` - Search destinations
- `POST /api/v1/trips/plan` - Plan a new trip
- `GET /api/v1/trips/{id}` - Get trip details
- `GET /api/v1/trips/user/{userId}` - Get user trips
- `POST /api/v1/ai/chat` - Chat with AI assistant

## 🏗️ Architecture

### Backend (`simple-backend.go`)
- Pure Go with standard library only
- No external dependencies
- In-memory data storage
- RESTful API with JSON responses
- CORS enabled for frontend access

### Frontend (`simple-frontend.html`)
- Single HTML file with embedded CSS and JavaScript
- No build process required
- Responsive design with CSS Grid
- Real-time chat interface
- Interactive trip planning form

## 🎯 Available Destinations

1. **Rishikesh** - Spiritual & Adventure
   - Vibes: serene, spiritual, adventure, mountains
   - Activities: River rafting, yoga, temples

2. **Goa** - Beach & Party
   - Vibes: beach, party, relaxing, nightlife
   - Activities: Beach activities, nightlife, water sports

3. **Shimla** - Hill Station
   - Vibes: mountains, serene, colonial, cool climate
   - Activities: Sightseeing, nature walks, colonial architecture

4. **Manali** - Adventure Hub
   - Vibes: mountains, adventure, snow, trekking
   - Activities: Paragliding, trekking, snow activities

## 💬 AI Assistant Features

The AI assistant can help with:
- Destination recommendations based on preferences
- Budget planning advice
- Activity suggestions
- Best time to visit information
- Vibe-based matching

### Sample Questions to Ask:
- "I want a peaceful mountain destination"
- "Suggest a budget-friendly adventure trip"
- "What's the best beach destination for partying?"
- "I love spiritual places, what do you recommend?"

## 🔧 Customization

### Adding New Destinations
Edit the `destinations` array in `simple-backend.go`:

```go
var destinations = []Destination{
    {
        ID:          "5",
        Name:        "Your Destination",
        City:        "City Name",
        State:       "State Name",
        Country:     "India",
        Vibes:       []string{"vibe1", "vibe2"},
        CostLevel:   "low/medium/high",
        Description: "Description here",
        MinDays:     3,
        MaxDays:     7,
    },
}
```

### Enhancing AI Responses
Modify the `generateAIResponse` function in `simple-backend.go` to add more intelligent responses based on keywords.

## 🌐 Network Configuration

The application uses:
- Backend: `localhost:8080`
- Frontend: Opens directly in browser (no server needed)
- CORS: Enabled for all origins in development

## 📱 Mobile Support

The application is fully responsive and works on:
- Desktop browsers
- Mobile browsers
- Tablets

## 🔒 Production Considerations

For production deployment:
1. Replace in-memory storage with a real database
2. Add authentication and user management
3. Implement rate limiting
4. Add proper error logging
5. Use environment variables for configuration
6. Add HTTPS support

## 🆚 Differences from Full Version

This simple version:
- ✅ No external dependencies
- ✅ Single file backend and frontend
- ✅ Instant startup
- ❌ No MongoDB (in-memory storage)
- ❌ No Ollama AI (rule-based responses)
- ❌ No Docker (direct execution)
- ❌ No complex frameworks

## 🛠️ Troubleshooting

### Backend won't start
- Check if Go is installed: `go version`
- Check if port 8080 is available
- Check for firewall blocking

### Frontend can't connect
- Ensure backend is running on localhost:8080
- Check browser console for CORS errors
- Try opening `simple-frontend.html` directly in browser

### Chat not working
- Verify backend API is responding: `http://localhost:8080/health`
- Check browser network tab for failed requests

## 📄 License

This project is open source and available under the MIT License.
