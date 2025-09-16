import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
    title: 'AuraTravel AI - Your Intelligent Travel Companion',
    description: 'Plan your perfect trip with AI-powered recommendations, personalized itineraries, and smart travel insights.',
    keywords: 'travel, AI, trip planning, recommendations, itinerary, vacation, tourism',
    authors: [{ name: 'AuraTravel Team' }],
    viewport: 'width=device-width, initial-scale=1',
}

export default function RootLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <html lang="en">
            <body className={inter.className}>
                <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
                    {children}
                </div>
            </body>
        </html>
    )
}