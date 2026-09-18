# AI Translator - 智能翻译工具

## 项目概述
AI驱动的智能翻译工具 - 支持多语言互译、术语管理、批量翻译

## 核心功能
- 40+语言互译（DeepSeek API）
- 实时翻译界面
- 历史记录管理
- 常用短语收藏
- 术语库管理
- 批量翻译
- 导出翻译结果（JSON/TXT/CSV）

## 目标用户
- 内容创作者（多语言内容发布）
- 跨境电商卖家
- 技术文档翻译
- 日常学习翻译

## 技术栈
- 后端：Go + Gin + SQLite
- 前端：Next.js 14 + Tailwind CSS
- AI翻译：DeepSeek API
- 部署：Vercel + Render

## API设计
- POST /api/translate - 翻译文本
- GET /api/history - 获取历史
- POST /api/terms - 添加术语
- GET /api/terms - 获取术语库
- DELETE /api/terms/:id - 删除术语
- GET /api/stats - 统计信息

## 商业模式
- 免费版：100次/天
- 专业版：$9.9/月 - 无限翻译
- 团队版：$29/月 - 共享术语库
