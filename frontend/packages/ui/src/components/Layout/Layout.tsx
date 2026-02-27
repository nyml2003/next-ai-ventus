import React from 'react';
import './style.css';

export interface HeaderProps {
  siteName: string;
  navLinks: Array<{ name: string; href: string }>;
  loginHref?: string;
}

export const Header: React.FC<HeaderProps> = ({ siteName, navLinks, loginHref }) => {
  return (
    <header className="header">
      <div className="header__container">
        <a href="/" className="header__logo">
          {siteName}
        </a>
        <nav className="header__nav">
          {navLinks.map((link) => (
            <a key={link.href} href={link.href} className="header__link">
              {link.name}
            </a>
          ))}
          {loginHref && (
            <a href={loginHref} className="header__link header__link--login">
              登录
            </a>
          )}
        </nav>
      </div>
    </header>
  );
};

export interface FooterProps {
  copyright: string;
}

export const Footer: React.FC<FooterProps> = ({ copyright }) => {
  return (
    <footer className="footer">
      <div className="footer__container">
        <p className="footer__copyright">{copyright}</p>
      </div>
    </footer>
  );
};

export const Container: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return <div className="container">{children}</div>;
};
