# Real-Time Todo App

A modern, real-time todo application built with Go and Next.js, featuring WebSocket connections for instant updates and Redis for data persistence.

## Tech Stack
- **Backend**: Go with Gorilla Mux
- **Frontend**: Next.js 14 with TypeScript and Tailwind CSS
- **Database**: Redis
- **Real-time**: WebSocket for live updates

## Features
- ✨ Real-time updates across all connected clients
- 🎨 Modern dark UI with gradient accents
- 🔄 Instant state synchronization
- 💾 Redis persistence
- ✅ Task completion tracking
- 📱 Responsive design

## Architecture
- RESTful API for CRUD operations
- WebSocket connection for real-time updates
- Redis Pub/Sub for event broadcasting
- Type-safe frontend with TypeScript

## Getting Started
1. Start Redis (Docker)
2. Run Go backend
3. Start Next.js frontend
4. Open `http://localhost:3000`
