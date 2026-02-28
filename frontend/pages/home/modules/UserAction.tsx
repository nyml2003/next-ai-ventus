import { useModuleData } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface UserActionData {
  isLoggedIn: boolean;
  loginHref: string;
  user?: {
    name: string;
    avatar?: string;
  };
}

interface UserActionProps {
  pageProps: PageProps;
}

export const UserAction: React.FC<UserActionProps> = () => {
  const { data, loading } = useModuleData<UserActionData>();
  
  if (loading) {
    return <span>加载中...</span>;
  }
  
  // 已登录状态
  if (data?.isLoggedIn && data.user) {
    return (
      <a href="/pages/admin-posts/index.html" style={{ textDecoration: 'none', color: '#333' }}>
        {data.user.name}
      </a>
    );
  }
  
  // 未登录状态
  return (
    <a href={data?.loginHref || "/pages/login/index.html"} style={{ textDecoration: 'none', color: '#333' }}>
      登录
    </a>
  );
};
