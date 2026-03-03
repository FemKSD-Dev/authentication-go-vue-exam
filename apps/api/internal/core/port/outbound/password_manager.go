package outbound

type PasswordManager interface {
	Encode(raw string) (string, error)
	Matches(encoded, raw string) error
}
