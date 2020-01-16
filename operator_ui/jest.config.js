module.exports = {
  preset: 'ts-jest/presets/js-with-ts',
  moduleDirectories: [
    'node_modules',
    '<rootDir>/src/',
    '<rootDir>/support/',
    '<rootDir>/__tests__',
  ],
  setupFiles: ['<rootDir>/jest.setup.js'],
  transformIgnorePatterns: ['node_modules/(?!(@chainlink/json-api-client)/)'],
  testPathIgnorePatterns: [
    '<rootDir>/dist/',
    '<rootDir>/tmp/',
    '<rootDir>/node_modules/',
    '<rootDir>/__tests__/.eslintrc.js',
  ],
  moduleNameMapper: {
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/__mocks__/fileMock.js',
    '\\.(css|less|sass|scss)$': '<rootDir>/__mocks__/styleMock.js',
    // handle declaration file having an enum value
    'core/store/models': '<rootDir>/@types/core/store/models.d.ts',
  },
}
