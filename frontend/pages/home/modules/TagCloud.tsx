import { useRequest } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface Tag {
  name: string;
  count: number;
}

interface TagCloudProps {
  pageProps: PageProps;
}

export const TagCloud: React.FC<TagCloudProps> = () => {
  const { data, loading } = useRequest<Tag[]>({
    scene: 'tag.list',
    params: {}
  });
  
  if (loading) {
    return <div style={{ padding: '16px' }}>加载中...</div>;
  }
  
  if (!data || data.length === 0) {
    return null;
  }
  
  return (
    <div style={{ padding: '16px', border: '1px solid #eee', borderRadius: '8px' }}>
      <h3 style={{ margin: '0 0 16px 0' }}>热门标签</h3>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
        {data.map((tag) => (
          <a
            key={tag.name}
            href={`/?tag=${tag.name}`}
            style={{
              padding: '4px 12px',
              background: '#f0f0f0',
              borderRadius: '4px',
              textDecoration: 'none',
              color: '#333',
              fontSize: '14px'
            }}
          >
            {tag.name} ({tag.count})
          </a>
        ))}
      </div>
    </div>
  );
};
