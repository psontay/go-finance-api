package util

const (
	ADMIN  = "ADMIN"
	CLIENT = "CLIENT"
)

func IsSupportedRole(role string) bool {
	switch role {
	case ADMIN, CLIENT:
		return true
	}
	return false
}
