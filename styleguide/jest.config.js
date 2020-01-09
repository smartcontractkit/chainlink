module.exports = {
  preset: 'ts-jest/presets/js-with-ts',
  moduleDirectories: ['node_modules', '<rootDir>/src/'],
  testPathIgnorePatterns: ['<rootDir>/dist/', '<rootDir>/tmp/'],
}
