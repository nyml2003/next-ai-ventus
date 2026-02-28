import { useRequest } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface UserInfo {
  name: string;
  avatar: string;
}

interface UserActionProps {
  pageProps: PageProps;
}

export const UserAction: React.FC<UserActionProps> = () => {
  const { data, loading } = useRequest<UserInfo>({
    scene: 'user.getCurrent',
    params: {}
  });
  
  if (loading) {
    return <span>加载中...</span>;
  }
  
  if (data) {
    return (
      <a href="/pages/admin-posts/index.html" style={{ textDecoration: 'none', color: '#333' }}>
        {data.name}
      </a>
    );
  }
  
  return (
    <a href="/pages/login/index.html" style={{ textDecoration: 'none', color: '#333' }}>
      登录
    </a>
  );
};
