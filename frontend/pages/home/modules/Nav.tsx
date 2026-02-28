import type { PageProps } from '@ventus/types';

interface NavProps {
  pageProps: PageProps;
}

export const Nav: React.FC<NavProps> = () => {
  return (
    <nav style={{ display: 'flex', gap: '24px' }}>
      <a href="/" style={{ textDecoration: 'none', color: '#333' }}>首页</a>
      <a href="/?tag=tech" style={{ textDecoration: 'none', color: '#333' }}>技术</a>
      <a href="/?tag=life" style={{ textDecoration: 'none', color: '#333' }}>生活</a>
      <a href="/pages/login/index.html" style={{ textDecoration: 'none', color: '#333' }}>登录</a>
    </nav>
  );
};
