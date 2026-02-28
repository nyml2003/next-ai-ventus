import { useModuleData } from '@ventus/store';
import type { PageProps } from '@ventus/types';

interface Tag {
  name: string;
  count: number;
  href: string;
}

interface TagCloudData {
  tags: Tag[];
}

interface TagCloudProps {
  pageProps: PageProps;
}

export const TagCloud: React.FC<TagCloudProps> = () => {
  const { data, loading } = useModuleData<TagCloudData>();
  
  if (loading) {
    return <div style={{ padding: '16px' }}>加载中...</div>;
  }
  
  const tags = data?.tags || [];
  
  if (tags.length === 0) {
    return null;
  }
  
  return (
    <div style={{ padding: '16px', background: '#f9f9f9', borderRadius: '8px' }}>
      <h3 style={{ margin: '0 0 12px 0', fontSize: '16px' }}>标签</h3>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
        {tags.map((tag) => (
          <a
            key={tag.name}
            href={tag.href}
            style={{
              padding: '4px 12px',
              background: '#fff',
              border: '1px solid #eee',
              borderRadius: '16px',
              textDecoration: 'none',
              color: '#666',
              fontSize: '14px',
            }}
          >
            {tag.name}
            <span style={{ marginLeft: '4px', color: '#999' }}>({tag.count})</span>
          </a>
        ))}
      </div>
    </div>
  );
};
