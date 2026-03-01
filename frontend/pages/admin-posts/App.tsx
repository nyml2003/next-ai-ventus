import React, { useEffect, useState } from 'react';
import { Button } from '@ventus/ui';
import { fetchPageData, adminAPI } from '@ventus/utils';
import './style.css';

interface MenuItem {
  name: string;
  icon: string;
  href: string;
  active: boolean;
}

interface AdminSidebarData {
  user: {
    name: string;
    avatar: string;
  };
  menu: MenuItem[];
}

interface AdminPostItem {
  id: string;
  title: string;
  slug: string;
  status: 'draft' | 'published';
  tags: string[];
  createdAt: string;
  updatedAt: string;
  href: string;
}

interface AdminPostListData {
  stats: {
    total: number;
    published: number;
    draft: number;
  };
  items: AdminPostItem[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
  newPostHref: string;
}

interface PageData {
  adminSidebar: AdminSidebarData;
  adminPostList: AdminPostListData;
}

export const App: React.FC = () => {
  const [data, setData] = useState<PageData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadData = async () => {
    try {
      const response = await fetchPageData({
        page: 'adminPosts',
        modules: ['adminSidebar', 'adminFilter', 'adminPostList'],
        params: { page: 1 },
      });

      setData({
        adminSidebar: response.modules.adminSidebar?.data as AdminSidebarData,
        adminPostList: response.modules.adminPostList?.data as AdminPostListData,
      });
    } catch (err) {
      if ((err as Error).message?.includes('unauthorized')) {
        window.location.href = '/pages/login/index.html';
        return;
      }
      setError(err instanceof Error ? err.message : '加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleDelete = async (id: string) => {
    if (!confirm('确定要删除这篇文章吗？')) return;

    try {
      await adminAPI('post.delete', { id });
      loadData();
    } catch (err) {
      alert('删除失败: ' + (err instanceof Error ? err.message : '未知错误'));
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="loading__spinner" />
        <p>加载中...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error">
        <p>加载失败: {error}</p>
        <button onClick={loadData}>重试</button>
      </div>
    );
  }

  if (!data) {
    return null;
  }

  const { adminSidebar, adminPostList } = data;

  return (
    <div className="admin-layout">
      <aside className="admin-sidebar">
        <div className="admin-sidebar__header">
          <h2>Ventus</h2>
        </div>
        <nav className="admin-sidebar__nav">
          {adminSidebar.menu.map((item) => (
            <a
              key={item.href}
              href={item.href}
              className={item.active ? 'active' : ''}
            >
              {item.name}
            </a>
          ))}
        </nav>
      </aside>

      <main className="admin-main">
        <header className="admin-main__header">
          <h1>文章管理</h1>
          <a href={adminPostList.newPostHref}>
            <Button variant="primary">新建文章</Button>
          </a>
        </header>

        <div className="stats">
          <div className="stat-card">
            <span className="stat-card__value">{adminPostList.stats.total}</span>
            <span className="stat-card__label">总文章</span>
          </div>
          <div className="stat-card">
            <span className="stat-card__value">{adminPostList.stats.published}</span>
            <span className="stat-card__label">已发布</span>
          </div>
          <div className="stat-card">
            <span className="stat-card__value">{adminPostList.stats.draft}</span>
            <span className="stat-card__label">草稿</span>
          </div>
        </div>

        <div className="post-table">
          <table>
            <thead>
              <tr>
                <th>标题</th>
                <th>状态</th>
                <th>标签</th>
                <th>更新时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {adminPostList.items.map((post) => (
                <tr key={post.id}>
                  <td>
                    <a href={post.href} className="post-title">
                      {post.title}
                    </a>
                  </td>
                  <td>
                    <span className={`status-badge ${post.status}`}>
                      {post.status === 'published' ? '已发布' : '草稿'}
                    </span>
                  </td>
                  <td>
                    {post.tags.map((tag) => (
                      <span key={tag} className="tag">
                        {tag}
                      </span>
                    ))}
                  </td>
                  <td>{post.updatedAt}</td>
                  <td>
                    <a href={post.href} className="action-link">
                      编辑
                    </a>
                    <button
                      className="action-link delete"
                      onClick={() => handleDelete(post.id)}
                    >
                      删除
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </div>
  );
};
