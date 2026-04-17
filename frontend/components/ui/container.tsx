import { cn } from "@/lib/utils";

interface ContainerProps {
  children: React.ReactNode;
  className?: string;
  size?: "sm" | "md" | "lg" | "xl" | "full";
}

const sizeMap = {
  sm: "max-w-2xl",
  md: "max-w-3xl",
  lg: "max-w-5xl",
  xl: "max-w-7xl",
  full: "max-w-[1400px]",
};

export function Container({ children, className, size = "md" }: ContainerProps) {
  return (
    <div className={cn("container mx-auto px-4 md:px-6 lg:px-8", sizeMap[size], className)}>
      {children}
    </div>
  );
}
