import { useRequest } from '@ventus/store';
import type { PageProps, PostListItem } from '@ventus/types';

interface PostListData {
  items: PostListItem[];
  pagination: {
    page: number;
    totalPages: number;
  };
}

interface PostListProps {
  pageProps: PageProps;
}

export const PostList: React.FC<PostListProps> = ({ pageProps }) => {
  const page = parseInt(pageProps.getQuery('page') || '1');
  const tag = pageProps.getQuery('tag');
  
  const { data, loading } = useRequest<PostListData>({
    scene: 'post.list',
    params: { page, pageSize: 10, tag }
  });
  
  if (loading) {
    return <div style={{ padding: '24px' }}>加载中...</div>;
  }
  
  if (!data || data.items.length === 0) {
    return (
      <div style={{ padding: '24px', textAlign: 'center' }}>
        <h3>还没有文章</h3>
        <p>登录后台创建你的第一篇文章吧</p>
      </div>
    );
  }
  
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
      {data.items.map((post) => (
        <article 
          key={post.id} 
          style={{ 
            padding: '16px', 
            border: '1px solid #eee', 
            borderRadius: '8px'
          }}
        >
          <a href={post.href} style={{ textDecoration: 'none', color: '#333' }}>
            <h2 style={{ margin: '0 0 8px 0' }}>{post.title}</h2>
          </a>
          <p style={{ margin: '0 0 8px 0', color: '#666' }}>{post.excerpt}</p>
          <div style={{ display: 'flex', gap: '8px', fontSize: '14px', color: '#999' }}>
            <span>{post.date}</span>
            {post.tags.map(tag => (
              <span key={tag} style={{ 
                padding: '2px 8px', 
                background: '#f0f0f0', 
                borderRadius: '4px' 
              }}>
                {tag}
              </span>
            ))}
          </div>
        </article>
      ))}
      
      {data.pagination.totalPages > 1 && (
        <div style={{ display: 'flex', gap: '8px', justifyContent: 'center', marginTop: '16px' }}>
          {Array.from({ length: data.pagination.totalPages }, (_, i) => i + 1).map((p) => (
            <a
              key={p}
              href={`/?page=${p}${tag ? `&tag=${tag}` : ''}`}
              style={{
                padding: '8px 12px',
                background: p === page ? '#333' : '#f0f0f0',
                color: p === page ? '#fff' : '#333',
                borderRadius: '4px',
                textDecoration: 'none'
              }}
            >
              {p}
            </a>
          ))}
        </div>
      )}
    </div>
  );
};
