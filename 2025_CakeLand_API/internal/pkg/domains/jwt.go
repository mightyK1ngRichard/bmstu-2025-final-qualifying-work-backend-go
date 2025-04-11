package domains

type JWTClaimsKeys string

const (
	KeyUserIDClaim JWTClaimsKeys = "userID"
	KeyEpxClaim    JWTClaimsKeys = "exp"
)

func (c JWTClaimsKeys) String() string {
	return string(c)
}
