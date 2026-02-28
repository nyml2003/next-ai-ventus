/**
 * @ventus/request - 错误码定义
 * 
 * 与后端 server/internal/interfaces/http/response/response.go 保持同步
 * 修改时需同步更新前后端
 */

/** 错误码枚举 - 与后端严格对齐 */
export enum ErrorCode {
  // ===== 通用成功 =====
  SUCCESS = 0,

  // ===== 通用错误 (1-99) =====
  /** 参数错误 */
  INVALID_PARAM = 1,
  /** 服务器内部错误 */
  INTERNAL_ERROR = 2,
  /** 未授权 */
  UNAUTHORIZED = 3,
  /** 禁止访问 */
  FORBIDDEN = 4,
  /** 资源不存在 */
  NOT_FOUND = 5,
  /** 方法不允许 */
  METHOD_NOT_ALLOWED = 6,
  /** 请求超时 */
  TIMEOUT = 7,

  // ===== 认证错误 (100-199) =====
  /** 认证失败 */
  AUTH_FAILED = 100,
  /** Token 无效 */
  TOKEN_INVALID = 101,
  /** Token 过期 */
  TOKEN_EXPIRED = 102,
  /** Token 缺失 */
  TOKEN_MISSING = 103,
  /** 用户名或密码错误 */
  INVALID_CREDENTIALS = 104,

  // ===== 文章错误 (200-299) =====
  /** 文章不存在 */
  POST_NOT_FOUND = 200,
  /** 文章已存在 */
  POST_ALREADY_EXISTS = 201,
  /** 链接标识已存在 */
  SLUG_EXISTS = 202,
  /** 标题无效 */
  INVALID_TITLE = 203,
  /** 内容无效 */
  INVALID_CONTENT = 204,
  /** 链接标识无效 */
  INVALID_SLUG = 205,
  /** 版本冲突（乐观锁） */
  VERSION_CONFLICT = 206,
  /** 状态无效 */
  INVALID_STATUS = 207,
  /** 标签无效 */
  INVALID_TAG = 208,

  // ===== BFF 模块错误 (300-399) =====
  /** 模块不存在 */
  MODULE_NOT_FOUND = 300,
  /** 模块执行错误 */
  MODULE_EXECUTE_ERROR = 301,

  // ===== 文件上传错误 (400-499) =====
  /** 上传失败 */
  UPLOAD_FAILED = 400,
  /** 文件类型无效 */
  INVALID_FILE_TYPE = 401,
  /** 文件过大 */
  FILE_TOO_LARGE = 402,
  /** 文件不存在 */
  FILE_NOT_FOUND = 403,
}

/** 中文错误消息映射 */
export const ErrorMessages: Record<ErrorCode, string> = {
  [ErrorCode.SUCCESS]: '成功',
  [ErrorCode.INVALID_PARAM]: '参数错误，请检查输入',
  [ErrorCode.INTERNAL_ERROR]: '服务器内部错误，请稍后重试',
  [ErrorCode.UNAUTHORIZED]: '未登录，请先登录',
  [ErrorCode.FORBIDDEN]: '无权限执行此操作',
  [ErrorCode.NOT_FOUND]: '请求的资源不存在',
  [ErrorCode.METHOD_NOT_ALLOWED]: '请求方法不允许',
  [ErrorCode.TIMEOUT]: '请求超时，请稍后重试',

  [ErrorCode.AUTH_FAILED]: '认证失败',
  [ErrorCode.TOKEN_INVALID]: '登录已失效，请重新登录',
  [ErrorCode.TOKEN_EXPIRED]: '登录已过期，请重新登录',
  [ErrorCode.TOKEN_MISSING]: '未提供认证信息',
  [ErrorCode.INVALID_CREDENTIALS]: '用户名或密码错误',

  [ErrorCode.POST_NOT_FOUND]: '文章不存在或已被删除',
  [ErrorCode.POST_ALREADY_EXISTS]: '文章已存在',
  [ErrorCode.SLUG_EXISTS]: '链接标识已被使用，请更换',
  [ErrorCode.INVALID_TITLE]: '标题不能为空',
  [ErrorCode.INVALID_CONTENT]: '内容不能为空',
  [ErrorCode.INVALID_SLUG]: '链接标识格式错误',
  [ErrorCode.VERSION_CONFLICT]: '文章已被其他人修改，请刷新后重试',
  [ErrorCode.INVALID_STATUS]: '文章状态无效',
  [ErrorCode.INVALID_TAG]: '标签格式错误',

  [ErrorCode.MODULE_NOT_FOUND]: '页面模块不存在',
  [ErrorCode.MODULE_EXECUTE_ERROR]: '页面模块加载失败',

  [ErrorCode.UPLOAD_FAILED]: '文件上传失败',
  [ErrorCode.INVALID_FILE_TYPE]: '不支持的文件类型',
  [ErrorCode.FILE_TOO_LARGE]: '文件大小超过限制',
  [ErrorCode.FILE_NOT_FOUND]: '文件不存在',
};

/** 英文错误消息映射（备用） */
export const ErrorMessagesEN: Record<ErrorCode, string> = {
  [ErrorCode.SUCCESS]: 'Success',
  [ErrorCode.INVALID_PARAM]: 'Invalid parameter',
  [ErrorCode.INTERNAL_ERROR]: 'Internal server error',
  [ErrorCode.UNAUTHORIZED]: 'Unauthorized',
  [ErrorCode.FORBIDDEN]: 'Forbidden',
  [ErrorCode.NOT_FOUND]: 'Resource not found',
  [ErrorCode.METHOD_NOT_ALLOWED]: 'Method not allowed',
  [ErrorCode.TIMEOUT]: 'Request timeout',

  [ErrorCode.AUTH_FAILED]: 'Authentication failed',
  [ErrorCode.TOKEN_INVALID]: 'Token invalid',
  [ErrorCode.TOKEN_EXPIRED]: 'Token expired',
  [ErrorCode.TOKEN_MISSING]: 'Token missing',
  [ErrorCode.INVALID_CREDENTIALS]: 'Invalid username or password',

  [ErrorCode.POST_NOT_FOUND]: 'Post not found',
  [ErrorCode.POST_ALREADY_EXISTS]: 'Post already exists',
  [ErrorCode.SLUG_EXISTS]: 'Slug already exists',
  [ErrorCode.INVALID_TITLE]: 'Invalid title',
  [ErrorCode.INVALID_CONTENT]: 'Invalid content',
  [ErrorCode.INVALID_SLUG]: 'Invalid slug',
  [ErrorCode.VERSION_CONFLICT]: 'Version conflict, please refresh and retry',
  [ErrorCode.INVALID_STATUS]: 'Invalid status',
  [ErrorCode.INVALID_TAG]: 'Invalid tag',

  [ErrorCode.MODULE_NOT_FOUND]: 'Module not found',
  [ErrorCode.MODULE_EXECUTE_ERROR]: 'Module execution error',

  [ErrorCode.UPLOAD_FAILED]: 'Upload failed',
  [ErrorCode.INVALID_FILE_TYPE]: 'Invalid file type',
  [ErrorCode.FILE_TOO_LARGE]: 'File too large',
  [ErrorCode.FILE_NOT_FOUND]: 'File not found',
};

/** 需要用户重试的错误码 */
export const RETRYABLE_ERROR_CODES = [
  ErrorCode.TIMEOUT,
  ErrorCode.INTERNAL_ERROR,
];

/** 需要用户刷新页面的错误码 */
export const REFRESHABLE_ERROR_CODES = [
  ErrorCode.VERSION_CONFLICT,
  ErrorCode.TOKEN_INVALID,
  ErrorCode.TOKEN_EXPIRED,
];

/** 需要重新登录的错误码 */
export const RELOGIN_ERROR_CODES = [
  ErrorCode.UNAUTHORIZED,
  ErrorCode.TOKEN_INVALID,
  ErrorCode.TOKEN_EXPIRED,
  ErrorCode.TOKEN_MISSING,
  ErrorCode.AUTH_FAILED,
];

/**
 * 获取错误消息
 * @param code - 错误码
 * @param locale - 语言，默认中文
 * @returns 错误消息
 */
export function getErrorMessage(code: number, locale: 'zh' | 'en' = 'zh'): string {
  const messages = locale === 'zh' ? ErrorMessages : ErrorMessagesEN;
  return messages[code as ErrorCode] || `未知错误 (code: ${code})`;
}

/**
 * 检查错误是否需要重试
 * @param code - 错误码
 */
export function isRetryableError(code: number): boolean {
  return RETRYABLE_ERROR_CODES.includes(code as ErrorCode);
}

/**
 * 检查错误是否需要刷新页面
 * @param code - 错误码
 */
export function isRefreshableError(code: number): boolean {
  return REFRESHABLE_ERROR_CODES.includes(code as ErrorCode);
}

/**
 * 检查错误是否需要重新登录
 * @param code - 错误码
 */
export function isReloginError(code: number): boolean {
  return RELOGIN_ERROR_CODES.includes(code as ErrorCode);
}

/**
 * 创建带错误码的 Error 对象
 */
export class APIError extends Error {
  code: ErrorCode;
  data?: unknown;

  constructor(code: ErrorCode, message?: string, data?: unknown) {
    super(message || getErrorMessage(code));
    this.code = code;
    this.data = data;
    this.name = 'APIError';
  }
}

export default ErrorCode;
