'use client'

import { useState } from 'react'
import Link from 'next/link'
import { motion } from 'framer-motion'
import {
    Mail,
    Lock,
    User,
    Eye,
    EyeOff,
    ArrowRight,
    Plane,
    Facebook,
    Apple
} from 'lucide-react'

export default function LoginPage() {
    const [isSignUp, setIsSignUp] = useState(false)
    const [showPassword, setShowPassword] = useState(false)
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        confirmPassword: '',
        firstName: '',
        lastName: '',
        agreeToTerms: false
    })

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value, type, checked } = e.target
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        // Handle authentication logic here
        console.log('Form submitted:', formData)
    }

    const handleSocialLogin = (provider: string) => {
        console.log(`Login with ${provider}`)
        // Handle social login logic here
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                {/* Header */}
                <motion.div
                    initial={{ opacity: 0, y: -20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6 }}
                    className="text-center"
                >
                    <Link href="/" className="inline-flex items-center space-x-2 mb-6">
                        <Plane className="w-8 h-8 text-primary-600" />
                        <span className="text-2xl font-bold text-gray-900">AuraTravel</span>
                        <span className="text-sm bg-primary-100 text-primary-800 px-2 py-1 rounded-full">AI</span>
                    </Link>
                    <h2 className="text-3xl font-bold text-gray-900">
                        {isSignUp ? 'Create your account' : 'Welcome back'}
                    </h2>
                    <p className="mt-2 text-gray-600">
                        {isSignUp
                            ? 'Start planning your perfect trip with AI'
                            : 'Sign in to continue your travel journey'
                        }
                    </p>
                </motion.div>

                {/* Main Form */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6, delay: 0.1 }}
                    className="bg-white rounded-2xl shadow-xl p-8 space-y-6"
                >
                    {/* Social Login Buttons */}
                    <div className="space-y-3">
                        <button
                            onClick={() => handleSocialLogin('Google')}
                            className="w-full flex items-center justify-center px-4 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                        >
                            <User className="w-5 h-5 mr-3" />
                            <span className="text-gray-700">Continue with Google</span>
                        </button>
                        <div className="grid grid-cols-2 gap-3">
                            <button
                                onClick={() => handleSocialLogin('Facebook')}
                                className="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                            >
                                <Facebook className="w-5 h-5 mr-2" />
                                <span className="text-gray-700">Facebook</span>
                            </button>
                            <button
                                onClick={() => handleSocialLogin('Apple')}
                                className="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                            >
                                <Apple className="w-5 h-5 mr-2" />
                                <span className="text-gray-700">Apple</span>
                            </button>
                        </div>
                    </div>

                    <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                            <div className="w-full border-t border-gray-300" />
                        </div>
                        <div className="relative flex justify-center text-sm">
                            <span className="px-2 bg-white text-gray-500">Or continue with email</span>
                        </div>
                    </div>

                    {/* Email/Password Form */}
                    <form onSubmit={handleSubmit} className="space-y-4">
                        {isSignUp && (
                            <div className="grid grid-cols-2 gap-4">
                                <div className="relative">
                                    <label htmlFor="firstName" className="sr-only">First Name</label>
                                    <div className="relative">
                                        <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                        <input
                                            id="firstName"
                                            name="firstName"
                                            type="text"
                                            required={isSignUp}
                                            value={formData.firstName}
                                            onChange={handleInputChange}
                                            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                            placeholder="First name"
                                        />
                                    </div>
                                </div>
                                <div className="relative">
                                    <label htmlFor="lastName" className="sr-only">Last Name</label>
                                    <div className="relative">
                                        <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                        <input
                                            id="lastName"
                                            name="lastName"
                                            type="text"
                                            required={isSignUp}
                                            value={formData.lastName}
                                            onChange={handleInputChange}
                                            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                            placeholder="Last name"
                                        />
                                    </div>
                                </div>
                            </div>
                        )}

                        <div className="relative">
                            <label htmlFor="email" className="sr-only">Email address</label>
                            <div className="relative">
                                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                <input
                                    id="email"
                                    name="email"
                                    type="email"
                                    autoComplete="email"
                                    required
                                    value={formData.email}
                                    onChange={handleInputChange}
                                    className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                    placeholder="Email address"
                                />
                            </div>
                        </div>

                        <div className="relative">
                            <label htmlFor="password" className="sr-only">Password</label>
                            <div className="relative">
                                <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                <input
                                    id="password"
                                    name="password"
                                    type={showPassword ? 'text' : 'password'}
                                    autoComplete={isSignUp ? 'new-password' : 'current-password'}
                                    required
                                    value={formData.password}
                                    onChange={handleInputChange}
                                    className="w-full pl-10 pr-12 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                    placeholder="Password"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                                >
                                    {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                                </button>
                            </div>
                        </div>

                        {isSignUp && (
                            <div className="relative">
                                <label htmlFor="confirmPassword" className="sr-only">Confirm Password</label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                                    <input
                                        id="confirmPassword"
                                        name="confirmPassword"
                                        type="password"
                                        autoComplete="new-password"
                                        required={isSignUp}
                                        value={formData.confirmPassword}
                                        onChange={handleInputChange}
                                        className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                        placeholder="Confirm password"
                                    />
                                </div>
                            </div>
                        )}

                        {isSignUp && (
                            <div className="flex items-center">
                                <input
                                    id="agreeToTerms"
                                    name="agreeToTerms"
                                    type="checkbox"
                                    required={isSignUp}
                                    checked={formData.agreeToTerms}
                                    onChange={handleInputChange}
                                    className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                                />
                                <label htmlFor="agreeToTerms" className="ml-2 block text-sm text-gray-700">
                                    I agree to the{' '}
                                    <Link href="/terms" className="text-primary-600 hover:text-primary-500">
                                        Terms of Service
                                    </Link>{' '}
                                    and{' '}
                                    <Link href="/privacy" className="text-primary-600 hover:text-primary-500">
                                        Privacy Policy
                                    </Link>
                                </label>
                            </div>
                        )}

                        <motion.button
                            whileHover={{ scale: 1.02 }}
                            whileTap={{ scale: 0.98 }}
                            type="submit"
                            className="w-full bg-primary-600 text-white py-3 px-4 rounded-lg hover:bg-primary-700 transition-colors font-semibold flex items-center justify-center gap-2"
                        >
                            {isSignUp ? 'Create Account' : 'Sign In'}
                            <ArrowRight className="w-4 h-4" />
                        </motion.button>
                    </form>

                    {!isSignUp && (
                        <div className="text-center">
                            <Link href="/forgot-password" className="text-sm text-primary-600 hover:text-primary-500">
                                Forgot your password?
                            </Link>
                        </div>
                    )}
                </motion.div>

                {/* Switch between Sign In/Sign Up */}
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.6, delay: 0.2 }}
                    className="text-center"
                >
                    <p className="text-gray-600">
                        {isSignUp ? 'Already have an account?' : "Don't have an account?"}{' '}
                        <button
                            onClick={() => setIsSignUp(!isSignUp)}
                            className="text-primary-600 hover:text-primary-500 font-semibold"
                        >
                            {isSignUp ? 'Sign in' : 'Sign up'}
                        </button>
                    </p>
                </motion.div>

                {/* Features Preview */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6, delay: 0.3 }}
                    className="bg-white/50 backdrop-blur-sm rounded-xl p-6"
                >
                    <h3 className="text-lg font-semibold text-gray-900 mb-3 text-center">
                        What you'll get with AuraTravel AI
                    </h3>
                    <ul className="space-y-2 text-sm text-gray-600">
                        <li className="flex items-center">
                            <div className="w-2 h-2 bg-primary-500 rounded-full mr-3"></div>
                            Personalized travel recommendations
                        </li>
                        <li className="flex items-center">
                            <div className="w-2 h-2 bg-primary-500 rounded-full mr-3"></div>
                            AI-powered itinerary planning
                        </li>
                        <li className="flex items-center">
                            <div className="w-2 h-2 bg-primary-500 rounded-full mr-3"></div>
                            Real-time travel insights
                        </li>
                        <li className="flex items-center">
                            <div className="w-2 h-2 bg-primary-500 rounded-full mr-3"></div>
                            Collaborative trip planning
                        </li>
                    </ul>
                </motion.div>
            </div>
        </div>
    )
}