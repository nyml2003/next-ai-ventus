import { marked } from 'marked';
import hljs from 'highlight.js';

// 配置 marked
marked.setOptions({
  breaks: true,
  gfm: true,
});

// 渲染 Markdown 为 HTML
export function renderMarkdown(content: string): string {
  return marked.parse(content, {
    async: false,
    renderer: {
      code(token) {
        const { text, lang } = token;
        if (lang && hljs.getLanguage(lang)) {
          try {
            const highlighted = hljs.highlight(text, { language: lang }).value;
            return `<pre><code class="hljs language-${lang}">${highlighted}</code></pre>`;
          } catch (e) {
            // fallback to plain text
          }
        }
        return `<pre><code class="hljs">${escapeHtml(text)}</code></pre>`;
      },
    },
  }) as string;
}

// HTML 转义
function escapeHtml(text: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  };
  return text.replace(/[&<>"']/g, (m) => map[m]);
}

// 提取目录
export interface TOCItem {
  level: number;
  text: string;
  anchor: string;
}

export function extractTOC(content: string): TOCItem[] {
  const tokens = marked.lexer(content);
  const toc: TOCItem[] = [];

  for (const token of tokens) {
    if (token.type === 'heading') {
      toc.push({
        level: token.depth,
        text: token.text,
        anchor: generateAnchor(token.text),
      });
    }
  }

  return toc;
}

// 生成锚点
function generateAnchor(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-');
}
