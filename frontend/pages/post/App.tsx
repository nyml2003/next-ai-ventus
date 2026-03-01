import React, { useEffect, useState } from 'react';
import { Header, Footer, Container } from '@ventus/ui';
import { fetchPageData, publicAPI, createPageProps } from '@ventus/utils';
import type { Post } from '@ventus/types';
import './style.css';

interface PostData {
  header: {
    siteName: string;
    navLinks: Array<{ name: string; href: string }>;
    loginHref: string;
  };
  article: Post;
  footer: {
    copyright: string;
  };
}

export const App: React.FC = () => {
  const [data, setData] = useState<PostData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 从 URL 参数获取 slug
  const pageProps = createPageProps();
  const slug = pageProps.getParam('slug') || '';

  useEffect(() => {
    const loadData = async () => {
      try {
        const response = await fetchPageData({
          page: 'post',
          modules: ['header', 'article', 'footer'],
          params: { slug },
        });

        const postData: PostData = {
          header: response.modules.header?.data as PostData['header'],
          article: response.modules.article?.data as Post,
          footer: response.modules.footer?.data as PostData['footer'],
        };

        setData(postData);

        // 记录阅读
        publicAPI('post.recordView', { id: postData.article.id }).catch(() => {
          // 忽略埋点错误
        });
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load post');
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [slug]);

  if (loading) {
    return (
      <div className="loading">
        <div className="loading__spinner" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="error">
        <p>文章不存在或已删除</p>
        <a href="/">返回首页</a>
      </div>
    );
  }

  const post = data.article;

  return (
    <div className="page">
      <Header
        siteName={data.header.siteName}
        navLinks={data.header.navLinks}
        loginHref={data.header.loginHref}
      />

      <main className="main">
        <Container>
          <article className="post-article">
            <header className="post-header">
              <h1 className="post-title">{post.title}</h1>
              <div className="post-meta">
                {post.publishedAt && (
                  <time>{post.publishedAt}</time>
                )}
                {post.wordCount && (
                  <span className="word-count">{post.wordCount} 字</span>
                )}
              </div>
              {post.tags.length > 0 && (
                <div className="post-tags">
                  {post.tags.map((tag) => (
                    <a key={tag} href={`/?tag=${tag}`} className="post-tag">
                      {tag}
                    </a>
                  ))}
                </div>
              )}
            </header>

            <div
              className="post-content"
              dangerouslySetInnerHTML={{ __html: post.html || '' }}
            />
          </article>
        </Container>
      </main>

      <Footer copyright={data.footer.copyright} />
    </div>
  );
};
