/** @type {import('next').NextConfig} */
const nextConfig = {
  webpack: (config) => {
    const rules = config.module.rules.find((rule) => typeof rule.oneOf === 'object').oneOf.filter((rule) => Array.isArray(rule.use));

    rules.forEach((rule) => {
      rule.use.forEach((moduleLoader) => {
        if (moduleLoader.loader !== undefined && moduleLoader.loader.includes('css-loader') && typeof moduleLoader.options.modules === 'object') {
          moduleLoader.options = {
            ...moduleLoader.options,
            modules: {
              ...moduleLoader.options.modules, // This is where we allow camelCase class names
              exportLocalsConvention: 'camelCase'
            }
          };
        }
      });
    });

    return config;
  },
  images: {
    minimumCacheTTL: 86400,
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost:8888',
      },
      {
        protocol: 'http',
        hostname: '127.0.0.1:8888',
      },
      {
        protocol: 'https',
        hostname: 'qiniu.einsition.com',
      },
      {
        protocol: 'https',
        hostname: 'img.einsition.com',
      },
    ],
  },
}

module.exports = nextConfig
