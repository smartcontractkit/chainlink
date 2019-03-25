module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  setupTestFrameworkScriptFile: '<rootDir>/jest.setup.ts',
  testPathIgnorePatterns: ['/node_modules/', '/client/']
}
