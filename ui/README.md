# Packet Size Manager UI

A modern Next.js application for managing and calculating optimal packet sizes.

## Features

- View current packet sizes
- Update packet sizes with comma-separated values
- Calculate optimal packet distribution for a given number of items
- Modern UI with animations and responsive design
- Error handling and connection status indicators

## Prerequisites

- Node.js 18.x or higher
- npm or yarn
- Go backend running on port 3000

## Getting Started

1. Install dependencies:

```bash
pnpm install
# or
yarn install
```

2. Start the development server:

```bash
pnpm run dev
# or
yarn dev
```

3. Open [http://localhost:3001](http://localhost:3001) in your browser.

## Important Notes

- This application requires the Go backend to be running on port 3000
- The Next.js application runs on port 3001 to avoid conflicts with the backend
- If the backend is not running, the UI will show a connection error and retry automatically

## Building for Production

```bash
pnpm run build
# or
yarn build
```

Then start the production server:

```bash
pnpm run start
# or
yarn start
``` 