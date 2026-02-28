import { useModuleData } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface LogoData {
  siteName: string;
  logo?: string;
  href: string;
}

interface LogoProps {
  pageProps: PageProps;
}

export const Logo: React.FC<LogoProps> = () => {
  const { data, loading } = useModuleData<LogoData>();
  
  if (loading) {
    return <span style={{ fontSize: '24px', fontWeight: 'bold' }}>Loading...</span>;
  }
  
  const siteName = data?.siteName || 'Ventus';
  const href = data?.href || '/';
  
  return (
    <a href={href} style={{ fontSize: '24px', fontWeight: 'bold', textDecoration: 'none', color: '#333' }}>
      {siteName}
    </a>
  );
};
