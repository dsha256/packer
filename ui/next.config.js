/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.API_URL || 'http://localhost:3000'}/api/:path*`,
      },
    ];
  },
  env: {
    PORT: process.env.PORT || 80,
  },
  experimental: {
    proxyTimeout: 60000, // 60 seconds
  },
  serverRuntimeConfig: {
    maxDuration: 60,
  },
};

module.exports = nextConfig; 