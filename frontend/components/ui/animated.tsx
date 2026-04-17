"use client"

import { motion, type HTMLMotionProps } from "framer-motion"
import { fadeIn, fadeInUp, scaleIn } from "@/lib/animations"

// 淡入容器
export function FadeIn({ 
  children, 
  ...props 
}: HTMLMotionProps<"div">) {
  return (
    <motion.div {...fadeIn} {...props}>
      {children}
    </motion.div>
  )
}

// 从下往上淡入
export function FadeInUp({ 
  children, 
  ...props 
}: HTMLMotionProps<"div">) {
  return (
    <motion.div {...fadeInUp} {...props}>
      {children}
    </motion.div>
  )
}

// 缩放淡入
export function ScaleIn({ 
  children, 
  ...props 
}: HTMLMotionProps<"div">) {
  return (
    <motion.div {...scaleIn} {...props}>
      {children}
    </motion.div>
  )
}

// 列表项动画
export function AnimatedListItem({ 
  children, 
  index = 0,
  ...props 
}: HTMLMotionProps<"div"> & { index?: number }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ 
        duration: 0.3, 
        delay: index * 0.05,
        ease: "easeOut" 
      }}
      {...props}
    >
      {children}
    </motion.div>
  )
}
