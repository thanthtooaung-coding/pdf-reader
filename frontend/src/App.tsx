import { Routes, Route, Navigate } from 'react-router-dom';
import AuthGuard, { PublicOnlyRoute } from '@/components/auth/AuthGuard';
import LoginPage from '@/pages/LoginPage';
import RegisterPage from '@/pages/RegisterPage';
import ReaderPage from '@/pages/ReaderPage';
import WorkspaceFilesPage from '@/pages/WorkspaceFilesPage';
import WorkspaceListPage from '@/pages/WorkspaceListPage';

export default function App() {
  return (
    <Routes>
      <Route element={<PublicOnlyRoute />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Route>

      <Route element={<AuthGuard />}>
        <Route path="/" element={<WorkspaceListPage />} />
        <Route path="/workspaces/:workspaceId" element={<WorkspaceFilesPage />} />
        <Route
          path="/workspaces/:workspaceId/files/:fileId"
          element={<ReaderPage />}
        />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
