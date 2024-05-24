module.exports = {
  semi: false,
  singleQuote: true,
  printWidth: 80,
  endOfLine: 'auto',
  tabWidth: 2,
  trailingComma: 'all',
  overrides: [
    {
      files: '*.sol',
      options: {
        parser: 'solidity-parse',
        printWidth: 120,
        tabWidth: 2,
        useTabs: false,
        singleQuote: false,
        bracketSpacing: false,
        explicitTypes: 'always',
      },
    },
  ],
}
