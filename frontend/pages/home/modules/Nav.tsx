import { useModuleData } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface NavData {
  links: Array<{ name: string; href: string }>;
}

interface NavProps {
  pageProps: PageProps;
}

export const Nav: React.FC<NavProps> = () => {
  const { data, loading } = useModuleData<NavData>();
  
  if (loading) {
    return <nav style={{ display: 'flex', gap: '24px' }}>加载中...</nav>;
  }
  
  const links = data?.links || [
    { name: '首页', href: '/' },
    { name: '技术', href: '/?tag=tech' },
    { name: '生活', href: '/?tag=life' },
  ];
  
  return (
    <nav style={{ display: 'flex', gap: '24px' }}>
      {links.map((link) => (
        <a key={link.name} href={link.href} style={{ textDecoration: 'none', color: '#333' }}>
          {link.name}
        </a>
      ))}
    </nav>
  );
};
