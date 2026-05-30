import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { loginRequest, loginVerify } from '@/api/auth';
import OTPForm from '@/components/auth/OTPForm';
import Button from '@/components/ui/Button';
import { useAuth } from '@/context/AuthProvider';
import { ApiError } from '@/lib/apiError';
import {
  LoginRequestSchema,
  type LoginRequest,
  type VerifyOTPRequest,
} from '@/schemas/apiSchema';

export default function LoginPage() {
  const navigate = useNavigate();
  const { setSession } = useAuth();
  const [step, setStep] = useState<'credentials' | 'otp'>('credentials');
  const [email, setEmail] = useState('');
  const [devOtp, setDevOtp] = useState<string | undefined>();
  const [formError, setFormError] = useState<string | null>(null);
  const [credentials, setCredentials] = useState<LoginRequest | null>(null);

  const loginMutation = useMutation({
    mutationFn: loginRequest,
    onSuccess: (data) => {
      setDevOtp(data.otp);
      setStep('otp');
      setFormError(null);
    },
    onError: (error: ApiError) => setFormError(error.message),
  });

  const verifyMutation = useMutation({
    mutationFn: loginVerify,
    onSuccess: (data) => {
      setSession(data.access_token, data.user);
      navigate('/', { replace: true });
    },
    onError: (error: ApiError) => setFormError(error.message),
  });

  const handleCredentialsSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const body = LoginRequestSchema.parse({
      email: formData.get('email'),
      username: formData.get('username'),
      password: formData.get('password'),
    });
    setEmail(body.email);
    setCredentials(body);
    loginMutation.mutate(body);
  };

  const handleOtpSubmit = (values: VerifyOTPRequest) => {
    setFormError(null);
    verifyMutation.mutate(values);
  };

  const handleResend = () => {
    if (!credentials) return;
    loginMutation.mutate(credentials);
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-6">
      <div className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
        <h1 className="text-xl font-semibold text-slate-900">Sign in</h1>
        <p className="mt-1 text-sm text-slate-500">
          DevDoc Companion — OTP login
        </p>

        <div className="mt-6">
          {step === 'credentials' ? (
            <form onSubmit={handleCredentialsSubmit} className="space-y-4">
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-slate-600">
                  Email
                </span>
                <input
                  name="email"
                  type="email"
                  required
                  className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
                />
              </label>
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-slate-600">
                  Username
                </span>
                <input
                  name="username"
                  required
                  className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
                />
              </label>
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-slate-600">
                  Password
                </span>
                <input
                  name="password"
                  type="password"
                  required
                  minLength={8}
                  className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
                />
              </label>
              {formError && (
                <p className="text-sm text-red-600">{formError}</p>
              )}
              <Button
                type="submit"
                className="w-full"
                disabled={loginMutation.isPending}
              >
                {loginMutation.isPending ? 'Sending OTP...' : 'Continue'}
              </Button>
            </form>
          ) : (
            <OTPForm
              email={email}
              devOtp={devOtp}
              onSubmit={handleOtpSubmit}
              onResend={handleResend}
              isSubmitting={verifyMutation.isPending}
              isResending={loginMutation.isPending}
              error={formError}
            />
          )}
        </div>

        <p className="mt-6 text-center text-sm text-slate-500">
          No account?{' '}
          <Link to="/register" className="font-medium text-slate-900 underline">
            Register
          </Link>
        </p>
      </div>
    </div>
  );
}
