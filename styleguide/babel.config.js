module.exports = {
  presets: [
    '@babel/preset-env',
    '@babel/preset-typescript',
    '@babel/preset-react',
  ],
  overrides: [
    {
      presets: [
        [
          '@babel/preset-env',
          { useBuiltIns: 'usage', corejs: { version: 3, proposals: true } },
        ],
      ],
    },
  ],
}
