import { useState, useEffect } from 'react';
import { Button, Input, TextArea } from '@ventus/ui';
import { fetchPageData, adminAPI, createPageProps } from '@ventus/utils';
import { renderMarkdown } from '@ventus/markdown';
import './style.css';

interface EditorData {
  id?: string;
  title: string;
  content: string;
  tags: string[];
  status: 'draft' | 'published';
  cover: string;
  version: number;
  isNew: boolean;
}

export const App: React.FC = () => {
  const [data, setData] = useState<EditorData | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [preview, setPreview] = useState('');
  const [tagInput, setTagInput] = useState('');

  const pageProps = createPageProps();
  const postId = pageProps.getParam('id');

  useEffect(() => {
    const loadData = async () => {
      try {
        const response = await fetchPageData({
          page: 'adminEditor',
          modules: ['editor', 'editorSettings'],
          params: { id: postId },
        });

        const editorData = response.modules.editor?.data as EditorData | undefined;
        if (editorData) {
          setData(editorData);
          setPreview(renderMarkdown(editorData.content || ''));
        }
      } catch (err) {
        if ((err as Error).message?.includes('unauthorized')) {
          window.location.href = '/pages/login/index.html';
          return;
        }
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [postId]);

  const handleContentChange = (value: string) => {
    if (!data) return;
    setData({ ...data, content: value });
    setPreview(renderMarkdown(value));
  };

  const handleAddTag = () => {
    if (!data || !tagInput.trim()) return;
    if (data.tags.includes(tagInput.trim())) {
      setTagInput('');
      return;
    }
    setData({ ...data, tags: [...data.tags, tagInput.trim()] });
    setTagInput('');
  };

  const handleRemoveTag = (tag: string) => {
    if (!data) return;
    setData({ ...data, tags: data.tags.filter((t: string) => t !== tag) });
  };

  const handleSave = async (publish: boolean = false) => {
    if (!data) return;

    setSaving(true);
    try {
      const sceneCode = data.isNew ? 'post.create' : 'post.update';
      const requestData: Record<string, unknown> = {
        title: data.title,
        content: data.content,
        tags: data.tags,
      };

      if (!data.isNew) {
        requestData.id = data.id;
        requestData.version = data.version;
      }

      if (publish) {
        requestData.status = 'published';
      }

      await adminAPI(sceneCode, requestData);
      window.location.href = '/pages/admin-posts/index.html';
    } catch (err) {
      alert('保存失败: ' + (err instanceof Error ? err.message : '未知错误'));
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="loading__spinner" />
      </div>
    );
  }

  if (!data) {
    return null;
  }

  return (
    <div className="editor-page">
      <header className="editor-header">
        <h1>{data.isNew ? '新建文章' : '编辑文章'}</h1>
        <div className="editor-header__actions">
          <Button variant="secondary" onClick={() => window.history.back()}>
            返回
          </Button>
          <Button
            variant="secondary"
            loading={saving}
            onClick={() => handleSave(false)}
          >
            保存草稿
          </Button>
          <Button
            variant="primary"
            loading={saving}
            onClick={() => handleSave(true)}
          >
            发布
          </Button>
        </div>
      </header>

      <div className="editor-layout">
        <div className="editor-main">
          <div className="editor-field">
            <label>标题</label>
            <Input
              value={data.title}
              onChange={(e) => setData({ ...data, title: e.target.value })}
              placeholder="请输入文章标题"
            />
          </div>

          <div className="editor-field">
            <label>内容 (Markdown)</label>
            <TextArea
              value={data.content}
              onChange={(e) => handleContentChange(e.target.value)}
              placeholder="请输入文章内容，支持 Markdown 语法"
              rows={20}
            />
          </div>
        </div>

        <div className="editor-sidebar">
          <div className="editor-panel">
            <h3>设置</h3>

            <div className="editor-field">
              <label>标签</label>
              <div className="tag-input">
                <Input
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
                  placeholder="输入标签后按回车"
                />
                <Button variant="secondary" size="sm" onClick={handleAddTag}>
                  添加
                </Button>
              </div>
              <div className="tag-list">
                {data.tags.map((tag) => (
                  <span key={tag} className="editor-tag">
                    {tag}
                    <button onClick={() => handleRemoveTag(tag)}>×</button>
                  </span>
                ))}
              </div>
            </div>

            <div className="editor-field">
              <label>状态</label>
              <div className="status-display">
                {data.status === 'published' ? '已发布' : '草稿'}
              </div>
            </div>
          </div>

          <div className="editor-panel">
            <h3>预览</h3>
            <div
              className="editor-preview"
              dangerouslySetInnerHTML={{ __html: preview }}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
