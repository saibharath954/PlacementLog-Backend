package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("SECRET"))

/*
Claims represents the JWT token claims structure.
Contains user ID, role, and standard JWT registered claims.
*/
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

/*
GenerateJwtToken creates a new JWT token with the specified user ID and role.
The token expires in 24 hours from creation.

Parameters:
- userID: The unique identifier of the user
- role: The role of the user ("user" or "admin")

Returns:
- string: The generated JWT token
- error: Any error that occurred during token generation
*/
func GenerateJwtToken(userID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/*
ValidateJwtToken validates a JWT token and extracts user information.
Verifies the token signature and expiration.

Parameters:
- tokenString: The JWT token to validate

Returns:
- string: The user ID from the token
- string: The role from the token
- error: Any error that occurred during validation
*/
func ValidateJwtToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Enforce HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, claims.Role, nil
	}

	return "", "", fmt.Errorf("invalid token claims")
}

/*
ValidateUserToken validates a JWT token and ensures it belongs to a user.
This function specifically checks that the token has a "user" role.

Parameters:
- tokenString: The JWT token to validate

Returns:
- string: The user ID from the token
- error: Any error that occurred during validation or if token is not a user token
*/
func ValidateUserToken(tokenString string) (string, error) {
	userID, role, err := ValidateJwtToken(tokenString)
	if err != nil {
		return "", err
	}
	if role != "user" {
		return "", fmt.Errorf("unauthorized: user token required")
	}
	return userID, nil
}

/*
ValidateAdminToken validates a JWT token and ensures it belongs to an admin.
This function specifically checks that the token has an "admin" role.

Parameters:
- tokenString: The JWT token to validate

Returns:
- string: The admin ID from the token
- error: Any error that occurred during validation or if token is not an admin token
*/
func ValidateAdminToken(tokenString string) (string, error) {
	userID, role, err := ValidateJwtToken(tokenString)
	if err != nil {
		return "", err
	}
	if role != "admin" {
		return "", fmt.Errorf("unauthorized: admin token required")
	}
	return userID, nil
}
