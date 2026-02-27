import React from 'react';
import './style.css';

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input: React.FC<InputProps> = ({
  label,
  error,
  className,
  ...props
}) => {
  const classes = ['input', error && 'input--error', className]
    .filter(Boolean)
    .join(' ');

  return (
    <div className="input-wrapper">
      {label && <label className="input__label">{label}</label>}
      <input className={classes} {...props} />
      {error && <span className="input__error">{error}</span>}
    </div>
  );
};

// TextArea 组件
export interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
}

export const TextArea: React.FC<TextAreaProps> = ({
  label,
  error,
  className,
  ...props
}) => {
  const classes = ['textarea', error && 'textarea--error', className]
    .filter(Boolean)
    .join(' ');

  return (
    <div className="input-wrapper">
      {label && <label className="input__label">{label}</label>}
      <textarea className={classes} {...props} />
      {error && <span className="input__error">{error}</span>}
    </div>
  );
};
