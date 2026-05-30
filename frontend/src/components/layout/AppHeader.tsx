import { Link } from 'react-router-dom';
import Button from '@/components/ui/Button';
import { useAuth } from '@/context/AuthProvider';

interface AppHeaderProps {
  title: string;
  breadcrumbs?: { label: string; to?: string }[];
}

export default function AppHeader({ title, breadcrumbs = [] }: AppHeaderProps) {
  const { user, logout } = useAuth();

  return (
    <header className="flex items-center justify-between border-b border-slate-200 bg-white px-6 py-4">
      <div>
        {breadcrumbs.length > 0 && (
          <nav className="mb-1 flex flex-wrap items-center gap-1 text-xs text-slate-500">
            {breadcrumbs.map((crumb, index) => (
              <span key={`${crumb.label}-${index}`} className="flex items-center gap-1">
                {index > 0 && <span>/</span>}
                {crumb.to ? (
                  <Link to={crumb.to} className="hover:text-slate-800">
                    {crumb.label}
                  </Link>
                ) : (
                  <span>{crumb.label}</span>
                )}
              </span>
            ))}
          </nav>
        )}
        <h1 className="text-lg font-semibold text-slate-900">{title}</h1>
      </div>
      <div className="flex items-center gap-3">
        {user && (
          <span className="text-sm text-slate-600">{user.username}</span>
        )}
        <Button variant="ghost" onClick={logout}>
          Logout
        </Button>
      </div>
    </header>
  );
}
