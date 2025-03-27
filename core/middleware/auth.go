package middleware

import (
	"net/http"
	"strings"
)

const (
	secretKey = "MIIETDCCAjSgAwIBAgIBETANBgkqhkiG9w0BAQsFADAYMRYwFAYDVQQDDA1KZXRQcm9maWxlIENBMB4XDTI0MDkyMDEyMTEyN1oXDTI2MDkyMjEyMTEyN1owHzEdMBsGA1UEAwwUcHJvZDJ5LWZyb20tMjAyNDA5MjAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7SH/XcUoMwkDi8JJPzXWWHWFdOZdrP2Dqkz2W8iUi650cwz2vdPEd0tMzosLAj7ifkFEHUyiuEcL//q9d9Op7ZsV23lpPXX8tFMLFwugoQ9D8jDLT/XP9pp/YukWkKF5jpNbaCvsVQkDdYkArBkYvhH3aN4v9BkEsXahfgLLOPe4IG2FDJNf9R4to9V1vt+m2UVJB0zV4a/sVMKUZLgqKmKKKOKoLrE3OjBlZlb+Q0z2N5dsW0hDEVRFGmBUAbHN/mp44MMMvEIFKfoLIGpgic92P2O6uFh75PI7mcultL6yuR48ajErx8CjjQEGOSnoq/8hD+yVE+6GW2gJa2CPvAgMBAAGjgZkwgZYwCQYDVR0TBAIwADAdBgNVHQ4EFgQUb5NERj05GyNerQ/Mjm9XH8HXtLIwSAYDVR0jBEEwP4AUo562SGdCEjZBvW3gubSgUouX8bOhHKQaMBgxFjAUBgNVBAMMDUpldFByb2ZpbGUgQ0GCCQDSbLGDsoN54TATBgNVHSUEDDAKBggrBgEFBQcDATALBgNVHQ8EBAMCBaAwDQYJKoZIhvcNAQELBQADggIBALq6VfVUjmPI3N/w0RYoPGFYUieCfRO0zVvD1VYHDWsN3F9buVsdudhxEsUb8t7qZPkDKTOB6DB+apgt2ZdKwok8S0pwifwLfjHAhO3b+LUQaz/VmKQW8gTOS5kTVcpM0BY7UPF8cRBqxMsdUfm5ejYk93lBRPBAqntznDY+DNc9aXOldFiACyutB1/AIh7ikUYPbpEIPZirPdAahroVvfp2tr4BHgCrk9z0dVi0tk8AHE5t7Vk4OOaQRJzy3lST4Vv6Mc0+0z8lNa+Sc3SVL8CrRtnTAs7YpD4fpI5AFDtchNrgFalX+BZ9GLu4FDsshVI4neqV5Jd5zwWPnwRuKLxsCO/PB6wiBKzdapQBG+P9z74dQ0junol+tqxd7vUV/MFsR3VwVMTndyapIS+fMoe+ZR5g+y44R8C7fXyVE/geg+JXQKvRwS0C5UpnS5FcGk+61b0e4U7pwO20RlwhEFHLSaP61p2TaVGo/TQtT/fWmrtV+HegAv9P3X3Se+xIVtJzQsk8QrB/w52IB3FKiAKl/KRn1egbMIs4uoNAkqNZ9Ih2P1NpiQnONFmkiAgeynJ+0FPykKdJQbV3Mx44jkaHIif4aFReTsYX1WUBNu/QerZRjn4FVSHRaZPSR5Oi82Wz0Nj7IY9ocTpLnXFrqkb/Kt3S6B9s2Kol3Lr1ElYA"
)

func AuthWithMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		// 接下来，你可以使用token进行进一步验证
		if token != secretKey {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func GetTaskSecretKey() string {
	return secretKey
}
