/**
 * 动画配置 - 简洁的渐变加载效果
 */

// 淡入动画
export const fadeIn = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 },
  transition: { duration: 0.3, ease: "easeOut" }
}

// 从下往上淡入
export const fadeInUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: 20 },
  transition: { duration: 0.4, ease: "easeOut" }
}

// 列表项渐变加载（带延迟）
export const staggerItem = (index: number) => ({
  initial: { opacity: 0, y: 10 },
  animate: { opacity: 1, y: 0 },
  transition: { 
    duration: 0.3, 
    delay: index * 0.05, // 每项延迟50ms
    ease: "easeOut" 
  }
})

// 容器动画（用于列表）
export const staggerContainer = {
  animate: {
    transition: {
      staggerChildren: 0.05
    }
  }
}

// 缩放淡入
export const scaleIn = {
  initial: { opacity: 0, scale: 0.95 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.95 },
  transition: { duration: 0.3, ease: "easeOut" }
}
