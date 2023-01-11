package keycloak

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	jose "gopkg.in/square/go-jose.v2"
)

type OpenIDAuth struct {
	config viper.Viper
}

func NewOpenIDAuth(config viper.Viper) *OpenIDAuth {
	return &OpenIDAuth{
		config: config,
	}
}

func (oidc *OpenIDAuth) Authenticate(r *http.Request, w http.ResponseWriter) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("jwt token not provided")
	}

	bearerToken := strings.Split(authHeader, "Bearer")
	if len(bearerToken) != 2 {
		return "", errors.New("authorization header does not follow the format 'Authorization: Bearer <string with JWT>'")
	}

	token := bearerToken[1]

	claims, err := oidc.parseIDToken(token)
	if err != nil {
		return "", err
	}

	audience := oidc.config.GetString("KEYCLOAK_OIDC_AUDIENCE")
	if err := oidc.verifyAudience(claims, audience); err != nil {
		return "", err
	}

	authorizedParty := oidc.config.GetString("KEYCLOAK_OIDC_AUTHORIZED_PARTY")
	if err := oidc.verifyAuthorizedParty(claims, authorizedParty); err != nil {
		return "", err
	}

	if err := oidc.verifyExpiry(claims); err != nil {
		return "", err
	}

	if err := oidc.getSignatureKeys(oidc.config); err != nil {
		return "", err
	}

	if err := oidc.verifySignature(token, oidc.config); err != nil {
		return "", err
	}

	issuer := oidc.config.GetString("KEYCLOAK_OIDC_ISSUER")
	if err := oidc.verifyIssuer(claims, issuer); err != nil {
		return "", err
	}

	userID, err := oidc.getTokenField(claims, "sub")
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (oidc *OpenIDAuth) parseIDToken(token string) (map[string]interface{}, error) {
	var claims map[string]interface{}

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 2 {
		return claims, errors.New("malformed jwt token")
	}

	parsedToken, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return claims, err
	}

	if err := json.Unmarshal(parsedToken, &claims); err != nil {
		return claims, err
	}

	return claims, nil
}

func (oidc *OpenIDAuth) verifyAudience(claims map[string]interface{}, audience string) error {
	var audiences []string

	tokenAud, ok := claims["aud"]
	if !ok {
		return errors.New("token issued without audience")
	}

	tokenAudiences, ok := tokenAud.([]interface{})
	if ok {
		for _, aud := range tokenAudiences {
			audiences = append(audiences, aud.(string))
		}
	} else {
		aud, ok := tokenAud.(string)
		if ok {
			audiences = append(audiences, aud)
		} else {
			return errors.New("invalid audience type")
		}
	}

	fail := false
	for _, aud := range audiences {
		if aud == audience {
			fail = true
		}
	}

	if !fail {
		return errors.New("token issued to another audience")
	}

	return nil
}

func (oidc *OpenIDAuth) verifyAuthorizedParty(claims map[string]interface{}, azp string) error {
	tokenAzp, _ := claims["azp"].(string)
	if len(tokenAzp) == 0 {
		return nil
	}

	if tokenAzp != azp {
		return errors.New("token issued to another authorized party")
	}

	return nil
}

func (oidc *OpenIDAuth) verifyExpiry(claims map[string]interface{}) error {
	tokenExp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("token issued without exp")
	}

	if time.Time(time.Unix(int64(tokenExp), 0)).Before(time.Now()) {
		return errors.New("token is expired")
	}

	return nil
}

func (oidc *OpenIDAuth) getSignatureKeys(config viper.Viper) error {
	var keySet jose.JSONWebKeySet

	keys := config.Get("oidc_keys")
	if keys != nil {
		return nil
	}

	client := &http.Client{}

	jwksURL := oidc.config.GetString("KEYCLOAK_OIDC_CERTS")
	req, err := http.NewRequest(http.MethodGet, jwksURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get jwt keys")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &keySet); err != nil {
		return err
	}

	config.Set("oidc_keys", keySet)

	return nil
}

func (oidc *OpenIDAuth) verifySignature(token string, config viper.Viper) error {
	jws, err := jose.ParseSigned(token)
	if err != nil {
		return err
	}

	keySet := config.Get("oidc_keys").(jose.JSONWebKeySet)

	fails := 0
	for _, key := range keySet.Keys {
		_, err = jws.Verify(&key)
		if err != nil {
			fails++
		}
	}
	if fails == len(keySet.Keys) {
		return errors.New("invalid token signature")
	}

	return nil
}

func (oidc *OpenIDAuth) verifyIssuer(claims map[string]interface{}, issuer string) error {
	tokenIssuer, ok := claims["iss"].(string)
	if !ok {
		return errors.New("token issuer is invalid")
	}

	if tokenIssuer != issuer {
		return errors.New("unrecognized token issuer")
	}

	return nil
}

func (oidc *OpenIDAuth) getTokenField(claims map[string]interface{}, field string) (string, error) {
	tokenField, ok := claims[field].(string)
	if !ok {
		return "", fmt.Errorf("token issued without %s", field)
	}

	return tokenField, nil
}
