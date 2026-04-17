"use client"

import Link from "next/link";
import { motion } from "framer-motion";
import { Home, TrendingUp, BookOpen, Users, Settings } from "lucide-react";
import { cn } from "@/lib/utils";

interface NavItemProps {
  href: string;
  icon: React.ReactNode;
  label: string;
  active?: boolean;
  index?: number;
}

function NavItem({ href, icon, label, active, index = 0 }: NavItemProps) {
  return (
    <motion.div
      initial={{ opacity: 0, x: -10 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ 
        duration: 0.3, 
        delay: index * 0.05,
        ease: "easeOut" 
      }}
    >
      <Link
        href={href}
        className={cn(
          "flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-colors",
          active
            ? "bg-muted text-foreground"
            : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
        )}
      >
        {icon}
        <span>{label}</span>
      </Link>
    </motion.div>
  );
}

export function LeftNav() {
  return (
    <nav className="sticky top-[5rem] space-y-1">
      <NavItem href="/" icon={<Home className="h-4 w-4" />} label="首页" active index={0} />
      <NavItem href="/hot" icon={<TrendingUp className="h-4 w-4" />} label="热榜" index={1} />
      <NavItem href="/topics" icon={<BookOpen className="h-4 w-4" />} label="话题" index={2} />
      <NavItem href="/users" icon={<Users className="h-4 w-4" />} label="用户" index={3} />
      
      <div className="pt-4 pb-2">
        <div className="px-4 text-xs font-semibold text-muted-foreground">个人中心</div>
      </div>
      
      <NavItem href="/settings" icon={<Settings className="h-4 w-4" />} label="设置" index={4} />
    </nav>
  );
}
