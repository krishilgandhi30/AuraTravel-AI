# AuraTravel AI Frontend

A modern, dynamic frontend for the AuraTravel AI platform built with Next.js 14, TypeScript, and Tailwind CSS.

## Features

- 🚀 **Modern Tech Stack**: Next.js 14, TypeScript, Tailwind CSS
- 🎨 **Beautiful UI**: Responsive design with smooth animations
- 🔥 **Firebase Integration**: Authentication and real-time database
- 🤖 **AI-Powered**: Smart trip planning and recommendations
- 📱 **Mobile-First**: Fully responsive across all devices
- ⚡ **Performance**: Optimized for speed and SEO

## Tech Stack

- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Authentication**: Firebase Auth
- **Database**: Firebase Firestore
- **Animations**: Framer Motion
- **Icons**: Lucide React
- **HTTP Client**: Axios

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Firebase project (for authentication and database)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/auratravel-ai.git
   cd auratravel-ai/frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   # or
   yarn install
   ```

3. **Set up environment variables**
   
   Copy the `.env.local` file and fill in your configuration:
   ```bash
   cp .env.local .env.local
   ```
   
   Update the following variables:
   ```env
   NEXT_PUBLIC_FIREBASE_API_KEY=your_firebase_api_key
   NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_project.firebaseapp.com
   NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_project_id
   NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your_project.appspot.com
   NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your_sender_id
   NEXT_PUBLIC_FIREBASE_APP_ID=your_app_id
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api
   ```

4. **Run the development server**
   ```bash
   npm run dev
   # or
   yarn dev
   ```

5. **Open your browser**
   Navigate to [http://localhost:3000](http://localhost:3000)

## Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── globals.css        # Global styles
│   │   ├── layout.tsx         # Root layout
│   │   ├── page.tsx          # Homepage
│   │   ├── login/            # Authentication pages
│   │   ├── plan/             # Trip planning
│   │   └── dashboard/        # User dashboard
│   ├── components/           # Reusable components
│   ├── lib/                  # Utility libraries
│   │   ├── firebase.ts       # Firebase configuration
│   │   └── api/              # API client
│   └── types/                # TypeScript type definitions
├── public/                   # Static assets
├── tailwind.config.js       # Tailwind CSS configuration
├── next.config.js           # Next.js configuration
└── package.json             # Dependencies and scripts
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

## Key Features

### 🏠 Homepage
- Hero section with animations
- Feature highlights
- Testimonials
- Call-to-action sections

### 🔐 Authentication
- Email/password login and registration
- Social authentication (Google, Facebook, Apple)
- Password reset functionality
- Protected routes

### 🗺️ Trip Planning
- Multi-step form with validation
- Destination search
- Preference selection
- Budget planning
- AI-powered recommendations

### 📊 Dashboard
- Trip management
- Statistics overview
- Search and filtering
- Trip status tracking

## Firebase Setup

1. **Create a Firebase project** at [Firebase Console](https://console.firebase.google.com)

2. **Enable Authentication**
   - Go to Authentication > Sign-in method
   - Enable Email/Password
   - Enable Google, Facebook, Apple (optional)

3. **Create Firestore Database**
   - Go to Firestore Database
   - Create database in test mode
   - Set up security rules

4. **Get configuration keys**
   - Go to Project Settings
   - Copy the config object
   - Update your `.env.local` file

## Deployment

### Vercel (Recommended)

1. **Connect your repository** to Vercel
2. **Set environment variables** in Vercel dashboard
3. **Deploy** - Vercel will automatically build and deploy

### Other Platforms

Build the project for production:
```bash
npm run build
npm run start
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `NEXT_PUBLIC_FIREBASE_API_KEY` | Firebase API key | Yes |
| `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` | Firebase auth domain | Yes |
| `NEXT_PUBLIC_FIREBASE_PROJECT_ID` | Firebase project ID | Yes |
| `NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET` | Firebase storage bucket | Yes |
| `NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID` | Firebase messaging sender ID | Yes |
| `NEXT_PUBLIC_FIREBASE_APP_ID` | Firebase app ID | Yes |
| `NEXT_PUBLIC_API_BASE_URL` | Backend API base URL | Yes |
| `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY` | Google Maps API key | No |
| `NEXT_PUBLIC_WEATHER_API_KEY` | Weather API key | No |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### Common Issues

1. **Module not found errors**
   ```bash
   npm install
   # or delete node_modules and reinstall
   rm -rf node_modules package-lock.json
   npm install
   ```

2. **Environment variables not working**
   - Ensure variables start with `NEXT_PUBLIC_`
   - Restart the development server
   - Check `.env.local` file location

3. **Firebase authentication errors**
   - Verify Firebase configuration
   - Check if authentication methods are enabled
   - Ensure domains are authorized in Firebase console

4. **Build errors**
   ```bash
   npm run lint
   npm run build
   ```

### Performance Optimization

- Images are optimized using Next.js Image component
- Code splitting with dynamic imports
- Tailwind CSS purging for smaller bundle size
- Framer Motion animations are optimized

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, email support@auratravel.ai or join our Discord community.

---

Built with ❤️ by the AuraTravel team