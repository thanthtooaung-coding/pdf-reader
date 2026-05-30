import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@/components/ui/Button';
import { VerifyOTPRequestSchema, type VerifyOTPRequest } from '@/schemas/apiSchema';

interface OTPFormProps {
  email: string;
  onSubmit: (values: VerifyOTPRequest) => void;
  onResend?: () => void;
  isSubmitting?: boolean;
  isResending?: boolean;
  devOtp?: string;
  error?: string | null;
}

export default function OTPForm({
  email,
  onSubmit,
  onResend,
  isSubmitting,
  isResending,
  devOtp,
  error,
}: OTPFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<VerifyOTPRequest>({
    resolver: zodResolver(VerifyOTPRequestSchema),
    defaultValues: { email, otp: devOtp ?? '' },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <p className="text-sm text-slate-600">
        Enter the 6-digit code sent to <strong>{email}</strong>
      </p>

      {devOtp && (
        <p className="rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
          Dev OTP: <code className="font-mono">{devOtp}</code>
        </p>
      )}

      <label className="block">
        <span className="mb-1 block text-xs font-medium text-slate-600">
          OTP code
        </span>
        <input
          {...register('otp')}
          inputMode="numeric"
          maxLength={6}
          className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm tracking-widest focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
        />
        {errors.otp && (
          <p className="mt-1 text-xs text-red-600">{errors.otp.message}</p>
        )}
      </label>

      {error && <p className="text-sm text-red-600">{error}</p>}

      <Button type="submit" disabled={isSubmitting} className="w-full">
        {isSubmitting ? 'Verifying...' : 'Verify'}
      </Button>

      {onResend && (
        <Button
          type="button"
          variant="ghost"
          className="w-full"
          disabled={isResending}
          onClick={onResend}
        >
          {isResending ? 'Resending...' : 'Resend code'}
        </Button>
      )}

      <input type="hidden" {...register('email')} value={email} readOnly />
    </form>
  );
}
