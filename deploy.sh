#!/bin/bash

# Subfinder Service Deployment Script
# This script helps deploy the Subfinder Service to various cloud environments

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default values
BACKEND_URL="http://localhost:8080"
ENVIRONMENT="local"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --backend-url)
      BACKEND_URL="$2"
      shift 2
      ;;
    --env)
      ENVIRONMENT="$2"
      shift 2
      ;;
    --help)
      echo "Usage: $0 [options]"
      echo "Options:"
      echo "  --backend-url URL   Set the backend URL (default: http://localhost:8080)"
      echo "  --env ENV           Set the deployment environment (local, aws, gcp, azure)"
      echo "  --help              Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

echo -e "${YELLOW}Deploying Subfinder Service to ${ENVIRONMENT} environment${NC}"
echo "========================================"

# Create .env file from .env.example
echo -e "\n${YELLOW}Creating .env file...${NC}"
cp .env.example .env
echo "BACKEND_URL=${BACKEND_URL}" >> .env
echo -e "${GREEN}Created .env file${NC}"

# Deploy based on environment
case $ENVIRONMENT in
  local)
    echo -e "\n${YELLOW}Deploying locally with Docker Compose...${NC}"
    docker compose down
    docker compose up -d
    
    if [ $? -ne 0 ]; then
      echo -e "${RED}Error: Failed to deploy locally${NC}"
      exit 1
    fi
    
    echo -e "${GREEN}Deployed successfully!${NC}"
    echo "Backend is available at: http://localhost:8080"
    echo "Frontend is available at: http://localhost:3000"
    ;;
    
  aws)
    echo -e "\n${YELLOW}Deploying to AWS...${NC}"
    echo "This would deploy to AWS ECS or EKS"
    # Add AWS deployment commands here
    echo -e "${GREEN}AWS deployment placeholder${NC}"
    ;;
    
  gcp)
    echo -e "\n${YELLOW}Deploying to Google Cloud...${NC}"
    echo "This would deploy to Google Cloud Run or GKE"
    # Add GCP deployment commands here
    echo -e "${GREEN}GCP deployment placeholder${NC}"
    ;;
    
  azure)
    echo -e "\n${YELLOW}Deploying to Azure...${NC}"
    echo "This would deploy to Azure Container Apps or AKS"
    # Add Azure deployment commands here
    echo -e "${GREEN}Azure deployment placeholder${NC}"
    ;;
    
  *)
    echo -e "${RED}Error: Unknown environment ${ENVIRONMENT}${NC}"
    exit 1
    ;;
esac

echo -e "\n${GREEN}Deployment completed!${NC}"
