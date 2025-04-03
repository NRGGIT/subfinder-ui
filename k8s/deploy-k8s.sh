#!/bin/bash

# Subfinder Service Kubernetes Deployment Script
# This script helps deploy the Subfinder Service to Kubernetes

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default values
BACKEND_IMAGE="subfinder-backend:latest"
FRONTEND_IMAGE="subfinder-frontend:latest"
SUBFINDER_HOST="subfinder.example.com"
NAMESPACE="default"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --backend-image)
      BACKEND_IMAGE="$2"
      shift 2
      ;;
    --frontend-image)
      FRONTEND_IMAGE="$2"
      shift 2
      ;;
    --host)
      SUBFINDER_HOST="$2"
      shift 2
      ;;
    --namespace)
      NAMESPACE="$2"
      shift 2
      ;;
    --help)
      echo "Usage: $0 [options]"
      echo "Options:"
      echo "  --backend-image IMAGE   Set the backend image (default: subfinder-backend:latest)"
      echo "  --frontend-image IMAGE  Set the frontend image (default: subfinder-frontend:latest)"
      echo "  --host HOST             Set the host for the ingress (default: subfinder.example.com)"
      echo "  --namespace NAMESPACE   Set the Kubernetes namespace (default: default)"
      echo "  --help                  Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

echo -e "${YELLOW}Deploying Subfinder Service to Kubernetes${NC}"
echo "========================================"
echo "Backend Image: ${BACKEND_IMAGE}"
echo "Frontend Image: ${FRONTEND_IMAGE}"
echo "Host: ${SUBFINDER_HOST}"
echo "Namespace: ${NAMESPACE}"
echo "========================================"

# Create namespace if it doesn't exist
echo -e "\n${YELLOW}Creating namespace ${NAMESPACE} if it doesn't exist...${NC}"
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

# Apply ConfigMaps
echo -e "\n${YELLOW}Applying ConfigMaps...${NC}"
kubectl apply -f backend/configmap.yaml -n ${NAMESPACE}
kubectl apply -f frontend/configmap.yaml -n ${NAMESPACE}

# Apply Services
echo -e "\n${YELLOW}Applying Services...${NC}"
kubectl apply -f backend/service.yaml -n ${NAMESPACE}
kubectl apply -f frontend/service.yaml -n ${NAMESPACE}

# Apply Deployments with image substitution
echo -e "\n${YELLOW}Applying Deployments...${NC}"
cat backend/deployment.yaml | sed "s|\${BACKEND_IMAGE:-subfinder-backend:latest}|${BACKEND_IMAGE}|g" | kubectl apply -f - -n ${NAMESPACE}
cat frontend/deployment.yaml | sed "s|\${FRONTEND_IMAGE:-subfinder-frontend:latest}|${FRONTEND_IMAGE}|g" | kubectl apply -f - -n ${NAMESPACE}

# Apply Ingress with host substitution
echo -e "\n${YELLOW}Applying Ingress...${NC}"
cat ../ingress.yaml | sed "s|\${SUBFINDER_HOST:-subfinder.example.com}|${SUBFINDER_HOST}|g" | kubectl apply -f - -n ${NAMESPACE}

echo -e "\n${GREEN}Deployment completed!${NC}"
echo "You can access the application at: http://${SUBFINDER_HOST}"
echo "Check the status with: kubectl get all -n ${NAMESPACE} -l app=subfinder"
