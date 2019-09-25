module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  testPathIgnorePatterns: [
    '/node_modules/',
    '<rootDir>/dist/',
    '<rootDir>/box/',
    '<rootDir>/v0.5/',
    '<rootDir>/build/',
    '<rootDir>/contracts/',
  ],
  modulePathIgnorePatterns: [
    '<rootDir>/dist/',
    '<rootDir>/box/',
    '<rootDir>/v0.5/',
    '<rootDir>/build/',
    '<rootDir>/contracts/',
  ],
  transformIgnorePatterns: [
    '/node_modules/',
    '<rootDir>/dist/',
    '<rootDir>/box/',
    '<rootDir>/v0.5/',
    '<rootDir>/build/',
    '<rootDir>/contracts/',
  ],
}
