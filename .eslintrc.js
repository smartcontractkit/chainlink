module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  env: {
    es6: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/eslint-recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:prettier/recommended',
    'prettier/@typescript-eslint',
  ],
  rules: {
    'object-shorthand': ['error', 'always'],
    'prettier/prettier': [
      'error',
      {},
      {
        usePrettierrc: true,
      },
    ],
    '@typescript-eslint/no-empty-function': 'off',
    '@typescript-eslint/no-unused-vars': 'error',
    '@typescript-eslint/no-empty-interface': 'off',
    '@typescript-eslint/explicit-function-return-type': 'off',
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/ban-ts-ignore': 'warn',
    '@typescript-eslint/no-non-null-assertion': 'error',
    '@typescript-eslint/no-use-before-define': [
      'error',
      { functions: false, typedefs: false },
    ],
  },

  overrides: [
    // enable jest for tests
    {
      files: ['**/*.test.ts', '**/*.test.js', '**/__mocks__/**/*.js'],
      env: {
        jest: true,
      },
    },
    {
      files: ['evm/v0.5/test/**/*', 'evm/box/**/*', 'evm/**/migrations/**/*'],
      env: { node: true, mocha: true },
      globals: {
        assert: 'readonly',
        artifacts: 'readonly',
        web3: 'readonly',
        contract: 'readonly',
      },
    },
    // add react linting for all of our react projects
    {
      files: [
        'explorer/client/**/*',
        'operator_ui/**/*',
        'feeds_ui/**/*',
        'styleguide/**/*',
        'tools/json-api-client/**/*',
        'tools/local-storage/**/*',
        'tools/redux/**/*',
      ],
      plugins: ['react-hooks'],
      extends: ['plugin:react/recommended'],
      env: {
        node: true,
        browser: true,
      },
      settings: {
        react: {
          version: 'detect',
        },
      },
      rules: {
        'react/prop-types': 'off',
        'react-hooks/rules-of-hooks': 'error',
        'react-hooks/exhaustive-deps': 'error',
      },
    },
  ],
}
