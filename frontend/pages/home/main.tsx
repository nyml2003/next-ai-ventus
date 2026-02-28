import React from "react";
import ReactDOM from "react-dom/client";
import { createOrchestration } from "@ventus/orchestration";
import { createRequest } from "@ventus/request";
import { homeConfig } from "./orchestration";
import { Logo, Nav, UserAction, PostList, TagCloud, Footer } from "./modules";

// 创建请求实例
const request = createRequest({ baseURL: "/api" });

// 创建编排系统
const orchestration = createOrchestration({
  config: homeConfig,
  request,
  modules: {
    Logo,
    Nav,
    UserAction,
    PostList,
    TagCloud,
    Footer,
  },
  resolver: (token) => {
    // 首页业务语义化 spacing
    const isMobile = window.innerWidth < 768;

    const map: Record<string, string> = {
      // 导航
      navPadding: isMobile ? "12px 16px" : "0 24px",
      navGap: isMobile ? "8px" : "16px",

      // 内容区
      pagePadding: isMobile ? "16px" : "24px",
      sectionGap: isMobile ? "16px" : "32px",
      contentGap: isMobile ? "12px" : "16px",
      contentPadding: isMobile ? "0" : "0",

      // 侧边栏
      sidebarGap: isMobile ? "12px" : "16px",
      sidebarPadding: isMobile ? "0" : "16px",

      // 页脚
      footerPadding: isMobile ? "24px 16px" : "32px 24px",
    };

    return () => map[token] || "0px";
  },
});

// 渲染
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <orchestration.Renderer />
  </React.StrictMode>,
);
