import { useModuleData } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface FooterData {
  copyright: string;
}

interface FooterProps {
  pageProps: PageProps;
}

export const Footer: React.FC<FooterProps> = () => {
  const { data, loading } = useModuleData<FooterData>();
  
  if (loading) {
    return <footer style={{ padding: '24px', textAlign: 'center' }}>加载中...</footer>;
  }
  
  const copyright = data?.copyright || `© ${new Date().getFullYear()} Ventus Blog`;
  
  return (
    <footer style={{ padding: '24px', textAlign: 'center', color: '#666', fontSize: '14px' }}>
      {copyright}
    </footer>
  );
};
