# AI Translator

## 功能特性
- 40+语言互译（DeepSeek API）
- 实时翻译界面
- 历史记录管理
- 常用短语收藏
- 术语库管理
- 批量翻译（待实现）
- 导出翻译结果（JSON/TXT/CSV）

## 技术栈
- 后端：Go + Gin + SQLite
- 前端：Next.js 14 + Tailwind CSS
- AI翻译：DeepSeek API

## 使用方法
1. 输入要翻译的文本
2. 选择源语言和目标语言
3. 点击"开始翻译"
4. 查看翻译结果并复制

## API端点
- `POST /api/translate` - 翻译文本
- `GET /api/languages` - 获取支持语言
- `GET /api/history` - 获取翻译历史
- `POST /api/terms` - 添加术语
- `GET /api/terms` - 获取术语库
- `DELETE /api/terms/:id` - 删除术语
- `GET /api/stats` - 获取统计信息
