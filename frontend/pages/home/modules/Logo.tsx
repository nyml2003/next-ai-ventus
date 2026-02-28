import type { PageProps } from '@ventus/types';

interface LogoProps {
  pageProps: PageProps;
}

export const Logo: React.FC<LogoProps> = () => {
  return (
    <a href="/" style={{ fontSize: '24px', fontWeight: 'bold', textDecoration: 'none', color: '#333' }}>
      Ventus
    </a>
  );
};
