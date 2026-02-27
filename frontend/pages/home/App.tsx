import React, { useEffect, useState } from "react";
import { Header, Footer, Container, PostCard } from "@ventus/ui";
import { fetchPageData } from "@ventus/utils";
import type { BFFPageResponse, PostListItem } from "@ventus/types";
import "./style.css";

interface HomeData {
  header: {
    siteName: string;
    logo: string;
    navLinks: Array<{ name: string; href: string }>;
    loginHref: string;
  };
  postList: {
    items: PostListItem[];
    pagination: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
  footer: {
    copyright: string;
  };
}

export const App: React.FC = () => {
  const [data, setData] = useState<HomeData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        const response = await fetchPageData({
          page: "home",
          modules: ["header", "postList", "footer"],
          params: { page: 1 },
        });

        console.log("BFF Response:", response);

        const homeData: HomeData = {
          header: response.modules.header?.data,
          postList: response.modules.postList?.data,
          footer: response.modules.footer?.data,
        };

        setData(homeData);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load data");
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

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
        <button onClick={() => window.location.reload()}>重试</button>
      </div>
    );
  }

  if (!data) {
    return null;
  }

  return (
    <div className="page">
      <Header
        siteName={data.header.siteName}
        navLinks={data.header.navLinks}
        loginHref={data.header.loginHref}
      />

      <main className="main">
        <Container>
          <div className="post-list">
            {data.postList.items.length === 0 ? (
              <div className="empty">
                <h2>还没有文章</h2>
                <p>登录后台创建你的第一篇文章吧</p>
                <a href={data.header.loginHref} className="btn">
                  前往管理后台
                </a>
              </div>
            ) : (
              <>
                {data.postList.items.map((post) => (
                  <PostCard key={post.id} post={post} />
                ))}

                {data.postList.pagination.totalPages > 1 && (
                  <div className="pagination">
                    {Array.from(
                      { length: data.postList.pagination.totalPages },
                      (_, i) => i + 1,
                    ).map((page) => (
                      <a
                        key={page}
                        href={`/?page=${page}`}
                        className={`pagination__item ${
                          page === data.postList.pagination.page ? "active" : ""
                        }`}
                      >
                        {page}
                      </a>
                    ))}
                  </div>
                )}
              </>
            )}
          </div>
        </Container>
      </main>

      <Footer copyright={data.footer.copyright} />
    </div>
  );
};
