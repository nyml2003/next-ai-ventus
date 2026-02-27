import React from 'react';
import './style.css';

export interface TagProps {
  name: string;
  href?: string;
  onClick?: () => void;
}

export const Tag: React.FC<TagProps> = ({ name, href, onClick }) => {
  const className = 'tag';

  if (href) {
    return (
      <a href={href} className={className} onClick={onClick}>
        {name}
      </a>
    );
  }

  return (
    <span className={className} onClick={onClick}>
      {name}
    </span>
  );
};
