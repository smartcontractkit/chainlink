module.exports = {
  preset: 'ts-jest/presets/js-with-ts',
  transform: {
    '.+\\.(css|styl|less|sass|scss|png|jpg|ttf|woff|woff2)$':
      'jest-transform-stub',
  },
  modulePaths: ['src'],
}
