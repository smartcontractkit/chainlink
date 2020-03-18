module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  testPathIgnorePatterns: ['/node_modules/', 'dist/'],
  testRunner: 'jest-circus/runner',
  testTimeout: 90000,
}
