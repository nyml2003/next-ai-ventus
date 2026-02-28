import '@testing-library/jest-dom';

// Mock window.location
Object.defineProperty(window, 'location', {
  writable: true,
  value: {
    href: 'http://localhost:3000/',
    pathname: '/',
    search: '',
  },
});
