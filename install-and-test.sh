#!/bin/bash

echo "=== Installing Go ==="
echo "Please run the following command to install Go:"
echo "sudo installer -pkg /tmp/go.pkg -target /"
echo ""
echo "After installation, run this script again to test the build."
echo ""

# Check if Go is installed
if command -v go &> /dev/null; then
    echo "✅ Go is installed: $(go version)"
    echo ""
    
    echo "=== Installing dependencies ==="
    go mod tidy
    
    echo "=== Building the application ==="
    make build
    
    if [ $? -eq 0 ]; then
        echo "✅ Build successful!"
        echo ""
        
        echo "=== Running tests ==="
        make test
        
        if [ $? -eq 0 ]; then
            echo "✅ Tests passed!"
            echo ""
            echo "=== Ready to run! ==="
            echo "Start the server with: make run"
            echo "Test the API with: ./test-examples.sh"
        else
            echo "❌ Tests failed"
        fi
    else
        echo "❌ Build failed"
    fi
else
    echo "❌ Go is not installed yet"
    echo "Please install Go first using the command above"
fi
