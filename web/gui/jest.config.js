module.exports = {
  moduleDirectories: [
    'node_modules',
    '<rootDir>/src/',
    '<rootDir>/support/'
  ],
  setupFiles: ['<rootDir>/jest.setup.js'],
  testPathIgnorePatterns: [
    '<rootDir>/dist/',
    '<rootDir>/tmp/',
    '<rootDir>/node_modules/'
  ]
}
