import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async redirects() {
    return [
      {
        source: '/purchase',
        destination: '/purchase',
        permanent: false,
      },
    ]
  },
};

export default nextConfig;
