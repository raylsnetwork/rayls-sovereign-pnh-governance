package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

// Helpers
var testJwtSecret = []byte("test-secret-key-for-unit-tests")

func newAuthServiceWithDefaults() (AuthService, *testutil.FakePrivateNetworkRepository) {
	repo := testutil.NewFakePrivateNetworkRepository()
	repo.NotFoundErr = ErrRecordNotFound // use core's sentinel so errors.Is matches
	svc := NewAuthService(repo, testJwtSecret, &testutil.StubLogger{})
	return svc, repo
}

func seedUser(repo *testutil.FakePrivateNetworkRepository, username, plainPassword string) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	repo.Networks[username] = &domain.PrivateNetwork{
		Username: username,
		Password: string(hashed),
	}
}

func generateValidToken(username string, secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func generateExpiredToken(username string, secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestAuthService_SignUp_StoresHashedPasswordInRepository(t *testing.T) {
	// A successful signup stores the user with a bcrypt-hashed password (not plaintext)
	svc, repo := newAuthServiceWithDefaults()

	err := svc.SignUp(context.Background(), "alice", "secret123")

	require.NoError(t, err)
	stored := repo.Networks["alice"]
	require.NotNil(t, stored)
	assert.Equal(t, "alice", stored.Username)
	// Password must be hashed, not plaintext
	assert.NotEqual(t, "secret123", stored.Password)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte("secret123")))
}

func TestAuthService_SignUp_DuplicateUsernameReturnsConflictError(t *testing.T) {
	// Attempting to sign up with an already-registered username returns a ConflictError
	svc, repo := newAuthServiceWithDefaults()
	repo.ConflictErr = ErrRecordConflict
	seedUser(repo, "alice", "original-password")

	err := svc.SignUp(context.Background(), "alice", "new-password")

	require.Error(t, err)
	var conflictErr *ConflictError
	assert.True(t, errors.As(err, &conflictErr))
	assert.Equal(t, "username", conflictErr.Resource)
}

func TestAuthService_Login_ReturnsValidJwtToken(t *testing.T) {
	// A successful login returns a JWT token that can be parsed with the same secret
	svc, repo := newAuthServiceWithDefaults()
	seedUser(repo, "alice", "correct-password")

	tokenString, err := svc.Login(context.Background(), "alice", "correct-password")

	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	// Token should be parseable and contain the correct subject
	parsed, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return testJwtSecret, nil
	})
	require.NoError(t, parseErr)
	claims := parsed.Claims.(jwt.MapClaims)
	assert.Equal(t, "alice", claims["sub"])
}

func TestAuthService_Login_UnknownUsernameReturnsGenericCredentialsError(t *testing.T) {
	// Login with non-existent user returns a generic "Invalid username or password" (not "user not found")
	svc, _ := newAuthServiceWithDefaults()

	_, err := svc.Login(context.Background(), "unknown-user", "any-password")

	require.Error(t, err)
	var valErr *ValidationError
	require.True(t, errors.As(err, &valErr))
	assert.Contains(t, valErr.Message, "Invalid username or password")
	assert.NotContains(t, err.Error(), "not found")
}

func TestAuthService_Login_WrongPasswordReturnsGenericCredentialsError(t *testing.T) {
	// Login with wrong password returns a generic "Invalid username or password"
	svc, repo := newAuthServiceWithDefaults()
	seedUser(repo, "alice", "correct-password")

	_, err := svc.Login(context.Background(), "alice", "wrong-password")

	require.Error(t, err)
	var valErr *ValidationError
	require.True(t, errors.As(err, &valErr))
	assert.Contains(t, valErr.Message, "Invalid username or password")
}

func TestAuthService_ValidateToken_ReturnsPrivateNetworkForValidToken(t *testing.T) {
	// A valid token resolves to the associated private network
	svc, repo := newAuthServiceWithDefaults()
	seedUser(repo, "alice", "password")
	tokenString := generateValidToken("alice", testJwtSecret)

	pn, err := svc.ValidateToken(context.Background(), tokenString)

	require.NoError(t, err)
	require.NotNil(t, pn)
	assert.Equal(t, "alice", pn.Username)
}

func TestAuthService_ValidateToken_DeletedUserReturnsValidationError(t *testing.T) {
	// A valid token for a user that no longer exists is rejected
	svc, _ := newAuthServiceWithDefaults()
	// Generate token for "alice" but don't seed her in the repo
	tokenString := generateValidToken("alice", testJwtSecret)

	pn, err := svc.ValidateToken(context.Background(), tokenString)

	require.Error(t, err)
	assert.Nil(t, pn)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestAuthService_SignUp_EmptyCredentialsReturnsValidationError(t *testing.T) {
	// SignUp with an empty username is rejected before touching the database
	svc, _ := newAuthServiceWithDefaults()

	err := svc.SignUp(context.Background(), "", "some-password")

	require.Error(t, err)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestAuthService_Login_EmptyCredentialsReturnsValidationError(t *testing.T) {
	// Login with an empty username is rejected before touching the database
	svc, _ := newAuthServiceWithDefaults()

	_, err := svc.Login(context.Background(), "", "some-password")

	require.Error(t, err)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestAuthService_Login_DatabaseErrorIsWrapped(t *testing.T) {
	// When the repository fails with a non-not-found error, Login returns an InternalError
	svc, repo := newAuthServiceWithDefaults()
	repo.Error = errors.New("connection reset by peer")

	_, err := svc.Login(context.Background(), "alice", "any-password")

	require.Error(t, err)
	var internalErr *InternalError
	assert.True(t, errors.As(err, &internalErr))
}

func TestAuthService_ValidateToken_ExpiredTokenReturnsValidationError(t *testing.T) {
	// A token whose expiry has passed is rejected even if the signature is otherwise valid
	svc, repo := newAuthServiceWithDefaults()
	seedUser(repo, "alice", "password")
	tokenString := generateExpiredToken("alice", testJwtSecret)

	pn, err := svc.ValidateToken(context.Background(), tokenString)

	require.Error(t, err)
	assert.Nil(t, pn)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestAuthService_ValidateToken_MalformedTokenReturnsValidationError(t *testing.T) {
	// A garbage string that is not a JWT is rejected immediately
	svc, _ := newAuthServiceWithDefaults()

	pn, err := svc.ValidateToken(context.Background(), "this.is.not.a.valid.jwt")

	require.Error(t, err)
	assert.Nil(t, pn)
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
}

func TestAuthService_ValidateToken_DatabaseErrorIsWrapped(t *testing.T) {
	// A structurally valid token but a failing repository returns an InternalError
	svc, repo := newAuthServiceWithDefaults()
	tokenString := generateValidToken("alice", testJwtSecret)
	repo.Error = errors.New("db connection timeout")

	pn, err := svc.ValidateToken(context.Background(), tokenString)

	require.Error(t, err)
	assert.Nil(t, pn)
	var internalErr *InternalError
	assert.True(t, errors.As(err, &internalErr))
}
