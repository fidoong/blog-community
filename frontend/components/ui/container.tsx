import { cn } from "@/lib/utils";

interface ContainerProps {
  children: React.ReactNode;
  className?: string;
  size?: "sm" | "md" | "lg" | "xl";
}

const sizeMap = {
  sm: "max-w-2xl",
  md: "max-w-3xl",
  lg: "max-w-5xl",
  xl: "max-w-7xl",
};

export function Container({ children, className, size = "md" }: ContainerProps) {
  return (
    <div className={cn("container mx-auto px-4", sizeMap[size], className)}>
      {children}
    </div>
  );
}
