package tlshelper

type (
	TlsClientOptions struct {
		UseTls             bool
		CertBase64         string
		KeyBase64          string
		RootCACertBase64   string
		CertFile           string
		KeyFile            string
		RootCACertFile     string
		InsecureSkipVerify bool
	}
)
