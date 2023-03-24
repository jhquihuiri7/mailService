package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	jwtRequest "github.com/golang-jwt/jwt/request"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (a *AuthRequest) ParseAuth(r io.Reader) {
	json.NewDecoder(r).Decode(a)
}
func (resp *AuthResponse) ParseAuth(r io.Reader) {
	json.NewDecoder(r).Decode(resp)
}
func (resp *AuthResponse) Marshal() string {
	JSONresponse, _ := json.Marshal(resp)
	return string(JSONresponse)
}

func (a *AuthRequest) GetUserCredentials(db *sql.DB) (int, []byte, error) {
	credentials, err := db.Query("SELECT id, password FROM users WHERE user = ?", a.User)
	var (
		id           int
		hashPassword string
	)
	if err != nil {
		log.Fatal(err)
	}
	for credentials.Next() {
		credentials.Scan(&id, &hashPassword)
	}

	if id == 0 {
		return 0, []byte{}, err
	} else {
		return id, []byte(hashPassword), nil
	}

}
func (a *AuthRequest) Login(db *sql.DB) AuthResponse {
	id, hashPassword, err := a.GetUserCredentials(db)
	if err != nil {
		return AuthResponse{Error: "Usuario no encontrado"}
	} else {
		err = bcrypt.CompareHashAndPassword(hashPassword, []byte(a.Password))
		if err != nil {
			return AuthResponse{Error: "Contraseña incorrecta"}
		}
		response := a.GenerateToken()
		response.Client = id
		response.User = a.User
		return response
	}
}
func InitKeys() {
	//Read Private and Public Key from directory
	privateBytes, _ := os.ReadFile("files/private.pem")
	publicBytes, _ := os.ReadFile("files/public.pem")
	//Parse Private and Public Key from []Byte

	PrivateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	PublicKey, _ = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
}
func (a *AuthRequest) GenerateToken() AuthResponse {
	claim := Claim{
		Credentials: a,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
			Issuer:    "Logiciel Applab",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	result, _ := token.SignedString(PrivateKey)
	return AuthResponse{Error: "", Token: result}
}
func (resp *AuthResponse) DecodeJWT() {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(resp.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	payload := claims["Credentials"]
	payloadMap := payload.(map[string]interface{})
	resp.User = fmt.Sprintf("%v", payloadMap["user"])
}
func (resp *AuthResponse) ValidateToken(req *http.Request, db *sql.DB) AuthResponse {
	var response AuthResponse
	token, err := jwtRequest.ParseFromRequest(req, jwtRequest.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return PublicKey, nil
	}, jwtRequest.WithClaims(&Claim{}))
	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				response.Error = "Id Expired"
			case jwt.ValidationErrorSignatureInvalid:
				response.Error = "No valid token"
			}
		default:
			response.Error = "No valid token"
		}
	}
	if token.Valid {
		response.Token = token.Raw
		response.DecodeJWT()
		client := AuthRequest{User: response.User}
		id, _, err := client.GetUserCredentials(db)
		if err != nil {
			response.Error = err.Error()
		} else {
			response.Client = id
		}
	}
	return response
}
