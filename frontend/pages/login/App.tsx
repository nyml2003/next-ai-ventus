import React, { useState } from 'react';
import { Button, Input } from '@ventus/ui';
import { login } from '@ventus/utils';
import './style.css';

export const App: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await login({ username, password });
      // 登录成功，跳转到管理后台
      window.location.href = '/pages/admin-posts/index.html';
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-box">
        <h1 className="login-box__title">Ventus Blog</h1>
        <p className="login-box__subtitle">管理后台登录</p>

        <form className="login-form" onSubmit={handleSubmit}>
          <Input
            label="用户名"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="请输入用户名"
            required
          />

          <Input
            label="密码"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="请输入密码"
            required
          />

          {error && <div className="login-form__error">{error}</div>}

          <Button
            type="submit"
            variant="primary"
            size="lg"
            loading={loading}
            style={{ width: '100%' }}
          >
            登录
          </Button>
        </form>

        <p className="login-box__hint">默认账号: admin / admin</p>
      </div>
    </div>
  );
};
