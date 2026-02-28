import { Header, Footer, Container } from "@ventus/ui";
import { PostList } from "./modules";
import "./style.css";

// 页面组件 - 数据由编排系统统一获取，模块通过 useModuleData 读取
export const App: React.FC = () => {
  return (
    <div className="page">
      <Header />
      
      <main className="main">
        <Container>
          <PostList />
        </Container>
      </main>
      
      <Footer />
    </div>
  );
};
