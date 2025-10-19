#!/bin/bash

# Voice Gateway Quick Start Script
# This script starts the full voice gateway stack

set -e

echo "🎙️  Voice Gateway - Quick Start"
echo "================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker not found. Will skip NATS container."
    USE_DOCKER=false
else
    USE_DOCKER=true
fi

# Build the binaries
echo "🔨 Building binaries..."
if [ ! -f "bin/gateway" ]; then
    echo "   Building gateway..."
    go build -o bin/gateway ./cmd/gateway
fi

echo "✅ Build complete"
echo ""

# Start NATS if Docker is available
if [ "$USE_DOCKER" = true ]; then
    echo "🚀 Starting NATS JetStream..."

    # Check if NATS is already running
    if docker ps | grep -q "nats"; then
        echo "   NATS already running"
    else
        docker run -d --name voice-gateway-nats \
            -p 4222:4222 \
            -p 8222:8222 \
            nats:latest -js

        echo "   NATS started on ports 4222 (client) and 8222 (monitoring)"
    fi
    echo ""
fi

# Start the gateway
echo "🎯 Starting Voice Gateway..."
echo "   Server will be available at: http://localhost:8080"
echo ""
echo "📝 To test:"
echo "   1. Open http://localhost:8080 in your browser"
echo "   2. Click 'Start Echo Test'"
echo "   3. Allow microphone access"
echo "   4. Speak and hear yourself back!"
echo ""
echo "Press Ctrl+C to stop"
echo ""
echo "================================"
echo ""

# Run the gateway
./bin/gateway
