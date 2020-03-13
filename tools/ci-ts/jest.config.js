module.exports = {
  preset: 'ts-jest',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  testEnvironment: 'node',
  testPathIgnorePatterns: ['/node_modules/', 'dist/'],
  testRunner: 'jest-circus/runner',
  testTimeout: 90000,
}
