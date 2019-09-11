const path = require('path')

export default pluginOptions => ({
  webpack: (currentWebpackConfig, state) => {
    return {
      ...currentWebpackConfig,
      resolve: {
        ...currentWebpackConfig.resolve,
        alias: {
          ...currentWebpackConfig.resolve.alias,
          'core/store/models': path.resolve(
            __dirname,
            '@types/core/store/models.d.ts',
          ),
          'core/store/presenters': path.resolve(
            __dirname,
            '@types/core/store/presenters.d.ts',
          ),
        },
      },
    }
  },
})
