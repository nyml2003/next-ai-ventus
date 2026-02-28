import type { PageProps } from '@ventus/types';

interface FooterProps {
  pageProps: PageProps;
}

export const Footer: React.FC<FooterProps> = () => {
  return (
    <footer style={{ textAlign: 'center', color: '#999', fontSize: '14px' }}>
      <p>&copy; 2024 Ventus Blog. All rights reserved.</p>
    </footer>
  );
};
