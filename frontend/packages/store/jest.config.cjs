/** @type {import('jest').Config} */
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src'],
  testMatch: ['**/__tests__/**/*.test.ts'],
  moduleNameMapper: {
    '^@ventus/types$': '<rootDir>/../types/src/index.ts',
    '^@ventus/request$': '<rootDir>/../request/src/index.ts'
  }
};
