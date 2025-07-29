#!/bin/bash

echo "Testing API Documentation Endpoints"
echo "=================================="

# Start server if not running
if ! curl -s http://localhost:8080/hello &> /dev/null; then
    echo "Starting server..."
    go run . &
    SERVER_PID=$!
    sleep 3
else
    echo "Server already running"
fi

echo -e "\n1. Testing OpenAPI YAML endpoint..."
if curl -s "http://localhost:8080/openapi.yaml" | head -3 | grep -q "openapi:"; then
    echo "✅ OpenAPI YAML is accessible"
else
    echo "❌ OpenAPI YAML endpoint failed"
fi

echo -e "\n2. Testing API docs HTML endpoint..."
if curl -s "http://localhost:8080/api-docs" | grep -q "swagger-ui"; then
    echo "✅ API documentation is accessible"
else  
    echo "❌ API documentation endpoint failed"
fi

echo -e "\n3. Checking OpenAPI spec validity..."
if curl -s "http://localhost:8080/openapi.yaml" | grep -q "Card Game API"; then
    echo "✅ OpenAPI spec contains expected title"
else
    echo "❌ OpenAPI spec validation failed"
fi

echo -e "\n4. Counting documented endpoints..."
ENDPOINT_COUNT=$(curl -s "http://localhost:8080/openapi.yaml" | grep -c "^  /")
echo "📊 Found $ENDPOINT_COUNT documented endpoints"

if [ $ENDPOINT_COUNT -gt 20 ]; then
    echo "✅ Good endpoint coverage"
else
    echo "⚠️  Low endpoint coverage"
fi

# Cleanup
if [ ! -z "$SERVER_PID" ]; then
    echo -e "\nStopping server..."
    kill $SERVER_PID 2>/dev/null
fi

echo -e "\n🎯 API Documentation Test Complete!"
echo "   Visit http://localhost:8080/api-docs for interactive documentation"