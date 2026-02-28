# 编排系统测试指南

## 运行测试

```bash
# 运行 orchestration 包的所有测试
cd frontend/packages/orchestration
pnpm test

# 带覆盖率
pnpm test --coverage

# 监视模式
pnpm test --watch

# 只运行特定文件
pnpm test BlockRenderer.test.tsx
```

## 测试类型

### 1. 单元测试 - 测试单个组件

```typescript
import { renderWithOrchestration } from '../test-utils';

test('Module 渲染', () => {
  const { container } = renderWithOrchestration({
    modules: { MyModule: MyComponent },
    config: {
      regions: [{
        id: 'content',
        type: 'content',
        block: {
          type: 'block',
          flexDirection: 'row',
          children: [{ type: 'module', name: 'MyModule' }]
        }
      }]
    }
  });
  
  expect(container.querySelector('.module-MyModule')).toBeInTheDocument();
});
```

### 2. 集成测试 - 测试数据流

```typescript
test('useRequest 获取数据', async () => {
  const { screen } = renderWithOrchestration({
    modules: { PostList },
    mockResponses: {
      'post.list:{}': { items: [{ id: '1', title: 'Test' }] }
    }
  });
  
  await waitFor(() => {
    expect(screen.getByText('Test')).toBeInTheDocument();
  });
});
```

### 3. 快照测试 - 验证结构

```typescript
test('页面结构快照', () => {
  const { container } = renderWithOrchestration({
    modules: { Header, Content },
    config: homeConfig
  });
  
  expect(container.firstChild).toMatchSnapshot();
});
```

## 最佳实践

1. **隔离测试**：每个测试独立运行，不要依赖其他测试的状态
2. **Mock API**：使用 `mockResponses` 模拟后端数据，不要真实请求
3. **测试关键路径**：重点测试数据流和交互，而不是样式细节
4. **异步处理**：数据获取是异步的，记得使用 `waitFor`

## 调试技巧

```typescript
// 打印组件结构
screen.debug();

// 打印特定元素
debug(container.querySelector('.block'));

// 查看 mock 调用
expect(mockRequest.call).toHaveBeenCalledWith({
  scene: 'post.list',
  params: { page: 1 }
});
```
