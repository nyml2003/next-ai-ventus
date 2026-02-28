module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  moduleNameMapper: {
    '^@ventus/types$': '<rootDir>/../types/src',
    '^@ventus/request$': '<rootDir>/../request/src',
    '^@ventus/store$': '<rootDir>/../store/src',
  },
  transform: {
    '^.+\\.tsx?$': ['ts-jest', {
      tsconfig: {
        jsx: 'react',
        esModuleInterop: true,
      },
    }],
  },
  testMatch: ['**/__tests__/**/*.test.ts(x)?'],
};
