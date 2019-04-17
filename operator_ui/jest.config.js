module.exports = {
  moduleDirectories: ['node_modules', '<rootDir>/src/', '<rootDir>/support/'],
  setupFiles: ['<rootDir>/jest.setup.js'],
  testPathIgnorePatterns: [
    '<rootDir>/dist/',
    '<rootDir>/tmp/',
    '<rootDir>/node_modules/'
  ],
  moduleNameMapper: {
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/__mocks__/fileMock.js',
    '\\.(css|less|sass|scss)$': '<rootDir>/__mocks__/styleMock.js'
  }
}
