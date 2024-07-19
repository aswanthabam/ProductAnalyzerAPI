/** @type {import('next').NextConfig} */
const nextConfig = {
  distDir: process.env.NEXT_PUBLIC_BUILD_PATH || "build",
  reactStrictMode: false,
  basePath: process.env.NEXT_PUBLIC_BASE_PATH || "/dashboard",
  output: "export",
  experimental: {
    missingSuspenseWithCSRBailout: false,
  },
};

module.exports = nextConfig;
