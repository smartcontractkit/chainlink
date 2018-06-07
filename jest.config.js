module.exports = {
  moduleDirectories: [
    'node_modules',
    '<rootDir>/gui/src/',
    '<rootDir>/gui/support/'
  ],
  setupFiles: ['<rootDir>/jest.setup.js'],
  testPathIgnorePatterns: [
    '<rootDir>/gui/dist/',
    '<rootDir>/gui/tmp/',
    '<rootDir>/node_modules/'
  ],
  moduleNameMapper: {
    '\\.(css|less|sass|scss)$': '<rootDir>/gui/__mocks__/styleMock.js'
  }
}
