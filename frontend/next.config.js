/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  experimental: {
    serverComponentsExternalPackages: ['@grpc/grpc-js', '@grpc/proto-loader'],
  },
};

module.exports = nextConfig;
