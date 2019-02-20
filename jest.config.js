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
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/gui/__mocks__/fileMock.js',
    '\\.(css|less|sass|scss)$': '<rootDir>/gui/__mocks__/styleMock.js'
  }
}
