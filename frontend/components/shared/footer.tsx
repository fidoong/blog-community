export function Footer() {
  return (
    <footer className="border-t bg-background">
      <div className="container mx-auto px-4 py-6 text-center text-sm text-muted-foreground">
        <p>© {new Date().getFullYear()} BlogCommunity. All rights reserved.</p>
      </div>
    </footer>
  );
}
