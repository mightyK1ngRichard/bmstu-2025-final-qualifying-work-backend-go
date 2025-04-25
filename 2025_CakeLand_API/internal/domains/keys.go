package domains

type JWTClaimsKeys string
type MetadataKey string

const (
	KeyUserIDClaim JWTClaimsKeys = "userID"
	KeyExpClaim    JWTClaimsKeys = "exp"

	KeyFingerprint   MetadataKey = "fingerprint"
	KeyAuthorization MetadataKey = "authorization"
)

func (c JWTClaimsKeys) String() string {
	return string(c)
}

func (k MetadataKey) String() string {
	return string(k)
}
