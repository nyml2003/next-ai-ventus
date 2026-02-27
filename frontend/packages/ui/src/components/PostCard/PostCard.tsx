import React from 'react';
import type { PostListItem } from '@ventus/types';
import { Tag } from '../Tag';
import './style.css';

export interface PostCardProps {
  post: PostListItem;
}

export const PostCard: React.FC<PostCardProps> = ({ post }) => {
  return (
    <article className="post-card">
      <a href={post.href} className="post-card__link">
        <h2 className="post-card__title">{post.title}</h2>
        <p className="post-card__excerpt">{post.excerpt}</p>
        <div className="post-card__meta">
          <time className="post-card__date">{post.date}</time>
          {post.tags.length > 0 && (
            <div className="post-card__tags">
              {post.tags.map((tag) => (
                <Tag key={tag} name={tag} />
              ))}
            </div>
          )}
        </div>
      </a>
    </article>
  );
};
