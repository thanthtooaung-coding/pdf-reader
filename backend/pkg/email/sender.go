package email

type Sender interface {
	SendOTP(to, otp string) error
	Enabled() bool
}
