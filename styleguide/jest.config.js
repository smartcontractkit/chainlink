module.exports = {
  preset: 'ts-jest',
  moduleDirectories: ['node_modules', '<rootDir>/src/'],
  testPathIgnorePatterns: ['<rootDir>/dist/', '<rootDir>/tmp/'],
}
