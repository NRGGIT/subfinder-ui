# Subfinder UI

A modern web interface for the Subfinder subdomain enumeration service, built with Nuxt.js and NuxtUI.

## Features

- **Job Submission**: Configure and submit subdomain enumeration jobs
- **Job Tracking**: Monitor the status of all jobs
- **Results Viewing**: Explore and search through subdomain results
- **Responsive Design**: Works on desktop and mobile devices

## Technology Stack

- **Frontend Framework**: [Nuxt.js](https://nuxt.com/) (Vue.js)
- **UI Components**: [NuxtUI](https://ui.nuxt.com/)
- **Styling**: Tailwind CSS
- **API Communication**: Nuxt's built-in fetch utilities

## Development

### Prerequisites

- Node.js 18+
- npm or yarn

### Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

### Building for Production

```bash
# Build the application
npm run build

# Start the production server
npm run start
```

## Docker

The UI is containerized and can be run using Docker:

```bash
# Build the Docker image
docker build -t subfinder-ui .

# Run the container
docker run -p 3000:3000 -e NUXT_PUBLIC_API_BASE_URL=http://localhost:8080 subfinder-ui
```

## Environment Variables

- `NUXT_PUBLIC_API_BASE_URL`: URL of the Subfinder API service (default: http://subfinder-service:8080)

## Using with Docker Compose

The UI is configured to work with the Subfinder service using Docker Compose:

```bash
# From the root directory (not the ui directory)
docker compose up -d
```

This will start both the Subfinder service and the UI.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
