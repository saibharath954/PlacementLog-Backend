#!/bin/bash

echo "Setting up Placement Log Backend..."

if [ ! -f .env ]; then
    echo "Creating .env file..."
    
    SECRET=$(openssl rand -base64 32)
    
    cat > .env << EOF

# JWT Secret Key (required for token generation)
SECRET=$SECRET

# Database Configuration (optional)
# DB_HOST=localhost
# DB_PORT=5432
# DB_NAME=placementlog
# DB_USER=postgres
# DB_PASSWORD=password

# Server Configuration
# PORT=8080

EOF
    
    echo ".env file created with generated SECRET"
    echo "Generated SECRET: $SECRET"
else
    echo ".env file already exists"
fi

# Check if SECRET is set
if grep -q "SECRET=" .env; then
    echo "SECRET environment variable is configured"
else
    echo "SECRET environment variable is missing from .env"
    echo "Please add SECRET=your-secret-key to the .env file"
fi

echo ""
echo "Next steps:"
echo "1. Run: make runServer"
echo "2. The server will start on http://localhost:8080"
echo "3. Make sure the frontend is configured to connect to this URL"
echo ""
echo "If you get 'token is malformed' errors, ensure the SECRET is set in .env" 