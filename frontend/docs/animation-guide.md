# 动画系统说明

## 概述

项目使用 Framer Motion 实现简洁的渐变加载动画，遵循"数据渐变加载"的设计理念。

## 动画配置

### 1. 文章卡片动画 (PostCard)

**位置**：`/components/features/post-card.tsx`

**动画参数**：
```typescript
initial={{ opacity: 0, y: 8 }}
animate={{ opacity: 1, y: 0 }}
transition={{ 
  duration: 0.2,        // 动画时长 200ms
  delay: index * 0.03,  // 每项延迟 30ms
  ease: "easeOut" 
}}
```

**优化策略**：
- 只对前 10 项添加动画（`index < 10`）
- 最大延迟限制为 0.3 秒
- 滚动加载的新内容（index >= 20）禁用动画，避免卡顿
- 移动距离从 10px 减少到 8px，更细腻

**使用方式**：
```tsx
<PostCard 
  post={post} 
  index={index} 
  disableAnimation={index >= 20}  // 禁用动画
/>
```

### 2. 侧边栏动画

**作者榜**：
```typescript
initial={{ opacity: 0, x: 10 }}
animate={{ opacity: 1, x: 0 }}
transition={{ 
  duration: 0.2,
  delay: index * 0.05,  // 每项延迟 50ms
  ease: "easeOut" 
}
```

**热门标签**：
```typescript
initial={{ opacity: 0, scale: 0.9 }}
animate={{ opacity: 1, scale: 1 }}
transition={{ 
  duration: 0.2,
  delay: index * 0.05,
  ease: "easeOut" 
}
```

### 3. Tab 切换动画

**下划线滑动**：
```tsx
<motion.div 
  layoutId="activeTab"
  transition={{ duration: 0.2, ease: "easeOut" }}
/>
```

使用 `layoutId` 实现流畅的共享布局动画。

### 4. 评论区动画

**位置**：`/components/shared/comment-section.tsx`

```typescript
initial={{ opacity: 0, y: 10 }}
animate={{ opacity: 1, y: 0 }}
transition={{ 
  duration: 0.3,
  delay: index * 0.05,
  ease: "easeOut" 
}
```

### 5. Header 动画

**位置**：`/components/shared/header.tsx`

```typescript
initial={{ opacity: 0, y: -10 }}
animate={{ opacity: 1, y: 0 }}
transition={{ duration: 0.3, ease: "easeOut" }}
```

## 性能优化原则

### 1. 限制动画数量
- 只对首屏可见内容添加动画
- 滚动加载的内容禁用动画
- 避免大量元素同时动画

### 2. 控制延迟时间
- 单项延迟：30-50ms
- 最大总延迟：300ms
- 避免延迟累加过大

### 3. 简化动画效果
- 优先使用 `opacity` 和 `transform`（GPU 加速）
- 避免动画 `width`、`height` 等触发重排的属性
- 移动距离控制在 10px 以内

### 4. 条件渲染
```typescript
// 根据条件决定是否启用动画
const shouldAnimate = !disableAnimation && index < 10;

<motion.div
  initial={shouldAnimate ? { opacity: 0 } : false}
  animate={shouldAnimate ? { opacity: 1 } : false}
/>
```

## 动画时长标准

| 场景 | 时长 | 说明 |
|------|------|------|
| 微交互 | 150-200ms | 按钮 hover、小元素淡入 |
| 内容加载 | 200-300ms | 卡片、列表项渐变 |
| 页面切换 | 300-400ms | 路由切换、大区块动画 |
| 布局动画 | 200-250ms | Tab 切换、共享元素 |

## 缓动函数

统一使用 `easeOut`：
- 快速启动，缓慢结束
- 符合自然物理规律
- 给人流畅、响应快的感觉

## 调试技巧

### 1. 临时禁用所有动画
```typescript
// 在组件中添加
const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

<motion.div
  initial={prefersReducedMotion ? false : { opacity: 0 }}
  // ...
/>
```

### 2. 调整动画速度
```typescript
// 开发时可以放慢动画，方便观察
transition={{ 
  duration: 0.2 * 3,  // 放慢 3 倍
  delay: index * 0.03 * 3 
}}
```

### 3. 查看动画性能
- Chrome DevTools > Performance
- 录制滚动操作
- 查看 FPS 和 GPU 使用情况

## 常见问题

### Q: 滚动时感觉有延迟？
A: 检查以下几点：
1. 是否对所有项都添加了动画？应该只对前 10-20 项
2. 延迟是否累加过大？限制最大延迟为 0.3s
3. 是否使用了触发重排的属性？改用 transform

### Q: 动画不流畅？
A: 
1. 确保使用 GPU 加速属性（opacity, transform）
2. 减少同时动画的元素数量
3. 检查是否有大量 DOM 操作

### Q: 如何完全禁用动画？
A:
```typescript
<PostCard disableAnimation={true} />
```

或在全局 CSS 中：
```css
* {
  animation-duration: 0s !important;
  transition-duration: 0s !important;
}
```
