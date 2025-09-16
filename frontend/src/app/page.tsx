'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { motion } from 'framer-motion'
import {
    MapPin,
    Compass,
    Star,
    Users,
    Calendar,
    ArrowRight,
    Plane,
    Camera,
    Heart,
    Globe,
    Sparkles
} from 'lucide-react'

export default function HomePage() {
    const [currentImageIndex, setCurrentImageIndex] = useState(0)

    const heroImages = [
        '/images/travel1.jpg',
        '/images/travel2.jpg',
        '/images/travel3.jpg'
    ]

    useEffect(() => {
        const interval = setInterval(() => {
            setCurrentImageIndex((prev) => (prev + 1) % heroImages.length)
        }, 5000)
        return () => clearInterval(interval)
    }, [heroImages.length])

    const features = [
        {
            icon: <Sparkles className="w-8 h-8" />,
            title: "AI-Powered Recommendations",
            description: "Get personalized travel suggestions based on your preferences and past trips"
        },
        {
            icon: <Compass className="w-8 h-8" />,
            title: "Smart Itinerary Planning",
            description: "Create optimized day-by-day itineraries with real-time adjustments"
        },
        {
            icon: <Globe className="w-8 h-8" />,
            title: "Global Destination Database",
            description: "Access comprehensive information about destinations worldwide"
        },
        {
            icon: <Users className="w-8 h-8" />,
            title: "Social Travel Planning",
            description: "Plan trips with friends and family with collaborative features"
        }
    ]

    const testimonials = [
        {
            name: "Sarah Johnson",
            location: "New York, USA",
            text: "AuraTravel AI helped me plan the perfect honeymoon in Bali. Every recommendation was spot on!",
            rating: 5
        },
        {
            name: "Mike Chen",
            location: "Toronto, Canada",
            text: "The AI suggestions saved me hours of research. My family trip to Europe was flawlessly planned.",
            rating: 5
        },
        {
            name: "Elena Rodriguez",
            location: "Madrid, Spain",
            text: "I love how it adapts to my budget and preferences. Best travel planning tool I've ever used!",
            rating: 5
        }
    ]

    return (
        <div className="min-h-screen">
            {/* Navigation */}
            <nav className="fixed top-0 w-full bg-white/90 backdrop-blur-md shadow-sm z-50">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        <div className="flex items-center space-x-2">
                            <Plane className="w-8 h-8 text-primary-600" />
                            <span className="text-2xl font-bold text-gray-900">AuraTravel</span>
                            <span className="text-sm bg-primary-100 text-primary-800 px-2 py-1 rounded-full">AI</span>
                        </div>
                        <div className="hidden md:flex items-center space-x-8">
                            <Link href="#features" className="text-gray-600 hover:text-primary-600 transition-colors">
                                Features
                            </Link>
                            <Link href="#how-it-works" className="text-gray-600 hover:text-primary-600 transition-colors">
                                How it Works
                            </Link>
                            <Link href="#testimonials" className="text-gray-600 hover:text-primary-600 transition-colors">
                                Testimonials
                            </Link>
                            <Link href="/login" className="bg-primary-600 text-white px-6 py-2 rounded-lg hover:bg-primary-700 transition-colors">
                                Get Started
                            </Link>
                        </div>
                    </div>
                </div>
            </nav>

            {/* Hero Section */}
            <section className="relative pt-16 pb-20 overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-r from-primary-600/90 to-secondary-600/90"></div>
                <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.8 }}
                        className="text-center"
                    >
                        <h1 className="text-5xl md:text-7xl font-bold text-white mb-6">
                            Travel Smarter with
                            <span className="block bg-gradient-to-r from-yellow-300 to-orange-300 bg-clip-text text-transparent">
                                AI-Powered Planning
                            </span>
                        </h1>
                        <p className="text-xl text-white/90 mb-8 max-w-3xl mx-auto">
                            Discover your perfect destination, create personalized itineraries, and explore the world
                            with intelligent recommendations tailored just for you.
                        </p>
                        <div className="flex flex-col sm:flex-row gap-4 justify-center">
                            <Link href="/login">
                                <motion.button
                                    whileHover={{ scale: 1.05 }}
                                    whileTap={{ scale: 0.95 }}
                                    className="bg-white text-primary-600 px-8 py-4 rounded-lg font-semibold text-lg hover:bg-gray-50 transition-colors flex items-center justify-center gap-2"
                                >
                                    Start Planning <ArrowRight className="w-5 h-5" />
                                </motion.button>
                            </Link>
                            <motion.button
                                whileHover={{ scale: 1.05 }}
                                whileTap={{ scale: 0.95 }}
                                className="border-2 border-white text-white px-8 py-4 rounded-lg font-semibold text-lg hover:bg-white hover:text-primary-600 transition-colors"
                            >
                                Watch Demo
                            </motion.button>
                        </div>
                    </motion.div>
                </div>

                {/* Floating Elements */}
                <motion.div
                    animate={{
                        y: [0, -20, 0],
                        rotate: [0, 5, 0]
                    }}
                    transition={{
                        duration: 4,
                        repeat: Infinity,
                        ease: "easeInOut"
                    }}
                    className="absolute top-1/4 left-10 text-white/20"
                >
                    <Camera className="w-12 h-12" />
                </motion.div>
                <motion.div
                    animate={{
                        y: [0, 15, 0],
                        rotate: [0, -5, 0]
                    }}
                    transition={{
                        duration: 3,
                        repeat: Infinity,
                        ease: "easeInOut",
                        delay: 1
                    }}
                    className="absolute top-1/3 right-10 text-white/20"
                >
                    <Heart className="w-10 h-10" />
                </motion.div>
            </section>

            {/* Features Section */}
            <section id="features" className="py-20 bg-white">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <motion.div
                        initial={{ opacity: 0 }}
                        whileInView={{ opacity: 1 }}
                        transition={{ duration: 0.8 }}
                        className="text-center mb-16"
                    >
                        <h2 className="text-4xl font-bold text-gray-900 mb-4">
                            Why Choose AuraTravel AI?
                        </h2>
                        <p className="text-xl text-gray-600 max-w-3xl mx-auto">
                            Our advanced AI technology combined with extensive travel data creates
                            the most personalized and efficient travel planning experience.
                        </p>
                    </motion.div>

                    <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
                        {features.map((feature, index) => (
                            <motion.div
                                key={index}
                                initial={{ opacity: 0, y: 20 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                transition={{ duration: 0.5, delay: index * 0.1 }}
                                className="text-center p-6 rounded-xl bg-gradient-to-br from-gray-50 to-gray-100 hover:shadow-lg transition-shadow"
                            >
                                <div className="text-primary-600 mb-4 flex justify-center">
                                    {feature.icon}
                                </div>
                                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                                    {feature.title}
                                </h3>
                                <p className="text-gray-600">
                                    {feature.description}
                                </p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* How it Works Section */}
            <section id="how-it-works" className="py-20 bg-gradient-to-br from-blue-50 to-indigo-100">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <motion.div
                        initial={{ opacity: 0 }}
                        whileInView={{ opacity: 1 }}
                        transition={{ duration: 0.8 }}
                        className="text-center mb-16"
                    >
                        <h2 className="text-4xl font-bold text-gray-900 mb-4">
                            How It Works
                        </h2>
                        <p className="text-xl text-gray-600 max-w-3xl mx-auto">
                            Three simple steps to your perfect trip
                        </p>
                    </motion.div>

                    <div className="grid md:grid-cols-3 gap-8">
                        {[
                            {
                                step: "01",
                                title: "Tell Us Your Preferences",
                                description: "Share your travel style, budget, interests, and dream destinations with our AI."
                            },
                            {
                                step: "02",
                                title: "Get AI Recommendations",
                                description: "Receive personalized suggestions for destinations, activities, and accommodations."
                            },
                            {
                                step: "03",
                                title: "Book & Explore",
                                description: "Finalize your itinerary and embark on your perfectly planned adventure."
                            }
                        ].map((item, index) => (
                            <motion.div
                                key={index}
                                initial={{ opacity: 0, x: index % 2 === 0 ? -20 : 20 }}
                                whileInView={{ opacity: 1, x: 0 }}
                                transition={{ duration: 0.6, delay: index * 0.2 }}
                                className="text-center"
                            >
                                <div className="bg-primary-600 text-white w-16 h-16 rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">
                                    {item.step}
                                </div>
                                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                                    {item.title}
                                </h3>
                                <p className="text-gray-600">
                                    {item.description}
                                </p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Testimonials Section */}
            <section id="testimonials" className="py-20 bg-white">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <motion.div
                        initial={{ opacity: 0 }}
                        whileInView={{ opacity: 1 }}
                        transition={{ duration: 0.8 }}
                        className="text-center mb-16"
                    >
                        <h2 className="text-4xl font-bold text-gray-900 mb-4">
                            What Our Travelers Say
                        </h2>
                        <p className="text-xl text-gray-600">
                            Join thousands of satisfied travelers worldwide
                        </p>
                    </motion.div>

                    <div className="grid md:grid-cols-3 gap-8">
                        {testimonials.map((testimonial, index) => (
                            <motion.div
                                key={index}
                                initial={{ opacity: 0, y: 20 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                transition={{ duration: 0.5, delay: index * 0.1 }}
                                className="bg-gray-50 p-6 rounded-xl"
                            >
                                <div className="flex items-center mb-4">
                                    {[...Array(testimonial.rating)].map((_, i) => (
                                        <Star key={i} className="w-5 h-5 text-yellow-400 fill-current" />
                                    ))}
                                </div>
                                <p className="text-gray-600 mb-4 italic">
                                    "{testimonial.text}"
                                </p>
                                <div>
                                    <div className="font-semibold text-gray-900">
                                        {testimonial.name}
                                    </div>
                                    <div className="text-sm text-gray-500">
                                        {testimonial.location}
                                    </div>
                                </div>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className="py-20 bg-gradient-to-r from-primary-600 to-secondary-600">
                <div className="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        whileInView={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.8 }}
                    >
                        <h2 className="text-4xl font-bold text-white mb-4">
                            Ready to Start Your Journey?
                        </h2>
                        <p className="text-xl text-white/90 mb-8">
                            Join millions of travelers who trust AuraTravel AI to plan their perfect trips.
                        </p>
                        <Link href="/login">
                            <motion.button
                                whileHover={{ scale: 1.05 }}
                                whileTap={{ scale: 0.95 }}
                                className="bg-white text-primary-600 px-8 py-4 rounded-lg font-semibold text-lg hover:bg-gray-50 transition-colors inline-flex items-center gap-2"
                            >
                                Get Started Free <ArrowRight className="w-5 h-5" />
                            </motion.button>
                        </Link>
                    </motion.div>
                </div>
            </section>

            {/* Footer */}
            <footer className="bg-gray-900 text-white py-12">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="grid md:grid-cols-4 gap-8">
                        <div>
                            <div className="flex items-center space-x-2 mb-4">
                                <Plane className="w-8 h-8 text-primary-400" />
                                <span className="text-2xl font-bold">AuraTravel</span>
                            </div>
                            <p className="text-gray-400">
                                Your intelligent travel companion for planning perfect trips with AI-powered recommendations.
                            </p>
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Product</h3>
                            <ul className="space-y-2 text-gray-400">
                                <li><Link href="#features" className="hover:text-white">Features</Link></li>
                                <li><Link href="#how-it-works" className="hover:text-white">How it Works</Link></li>
                                <li><Link href="/plan" className="hover:text-white">Plan Trip</Link></li>
                                <li><Link href="/dashboard" className="hover:text-white">Dashboard</Link></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Support</h3>
                            <ul className="space-y-2 text-gray-400">
                                <li><Link href="#" className="hover:text-white">Help Center</Link></li>
                                <li><Link href="#" className="hover:text-white">Contact Us</Link></li>
                                <li><Link href="#" className="hover:text-white">Privacy Policy</Link></li>
                                <li><Link href="#" className="hover:text-white">Terms of Service</Link></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Connect</h3>
                            <ul className="space-y-2 text-gray-400">
                                <li><Link href="#" className="hover:text-white">Twitter</Link></li>
                                <li><Link href="#" className="hover:text-white">Facebook</Link></li>
                                <li><Link href="#" className="hover:text-white">Instagram</Link></li>
                                <li><Link href="#" className="hover:text-white">LinkedIn</Link></li>
                            </ul>
                        </div>
                    </div>
                    <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
                        <p>&copy; 2024 AuraTravel AI. All rights reserved.</p>
                    </div>
                </div>
            </footer>
        </div>
    )
}