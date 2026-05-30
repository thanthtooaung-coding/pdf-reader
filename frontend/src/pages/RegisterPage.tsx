import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { registerRequest, registerResend, registerVerify } from '@/api/auth';
import OTPForm from '@/components/auth/OTPForm';
import Button from '@/components/ui/Button';
import { useAuth } from '@/context/AuthProvider';
import { ApiError } from '@/lib/apiError';
import {
  RegisterRequestSchema,
  type VerifyOTPRequest,
} from '@/schemas/apiSchema';

export default function RegisterPage() {
  const navigate = useNavigate();
  const { setSession } = useAuth();
  const [step, setStep] = useState<'details' | 'otp'>('details');
  const [email, setEmail] = useState('');
  const [devOtp, setDevOtp] = useState<string | undefined>();
  const [formError, setFormError] = useState<string | null>(null);

  const registerMutation = useMutation({
    mutationFn: registerRequest,
    onSuccess: (data) => {
      setDevOtp(data.otp);
      setStep('otp');
      setFormError(null);
    },
    onError: (error: ApiError) => setFormError(error.message),
  });

  const verifyMutation = useMutation({
    mutationFn: registerVerify,
    onSuccess: (data) => {
      setSession(data.access_token, data.user);
      navigate('/', { replace: true });
    },
    onError: (error: ApiError) => setFormError(error.message),
  });

  const resendMutation = useMutation({
    mutationFn: registerResend,
    onSuccess: (data) => {
      setDevOtp(data.otp);
      setFormError(null);
    },
    onError: (error: ApiError) => setFormError(error.message),
  });

  const handleDetailsSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const body = RegisterRequestSchema.parse({
      fullname: formData.get('fullname'),
      username: formData.get('username'),
      email: formData.get('email'),
      password: formData.get('password'),
    });
    setEmail(body.email);
    registerMutation.mutate(body);
  };

  const handleOtpSubmit = (values: VerifyOTPRequest) => {
    setFormError(null);
    verifyMutation.mutate(values);
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-6">
      <div className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
        <h1 className="text-xl font-semibold text-slate-900">Create account</h1>
        <p className="mt-1 text-sm text-slate-500">
          Register with email verification
        </p>

        <div className="mt-6">
          {step === 'details' ? (
            <form onSubmit={handleDetailsSubmit} className="space-y-4">
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-slate-600">
                  Full name
                </span>
                <input
                  name="fullname"
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
                  minLength={2}
                  className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
                />
              </label>
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
                disabled={registerMutation.isPending}
              >
                {registerMutation.isPending ? 'Sending OTP...' : 'Continue'}
              </Button>
            </form>
          ) : (
            <OTPForm
              email={email}
              devOtp={devOtp}
              onSubmit={handleOtpSubmit}
              onResend={() => resendMutation.mutate({ email })}
              isSubmitting={verifyMutation.isPending}
              isResending={resendMutation.isPending}
              error={formError}
            />
          )}
        </div>

        <p className="mt-6 text-center text-sm text-slate-500">
          Already have an account?{' '}
          <Link to="/login" className="font-medium text-slate-900 underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}
