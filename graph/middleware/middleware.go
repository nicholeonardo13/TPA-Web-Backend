package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nicholeonardo13/gqlgen-todos/graph/model"
	"net/http"
)

var userCtxKey = &model.ContextKey{"user"}
var userCtxWrite = &userCtxWrites{"write"}

type userCtxWrites struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie , err := r.Cookie("CookieJar" )

			// Allow unauthenticated users in
			//if header == "" {
			//	fmt.Println("GAK AUTH")
			//	next.ServeHTTP(w, r)
			//	return
			//}

			writes := context.WithValue(r.Context(), userCtxWrite, &w)

			// and call the next with our new context
			r = r.WithContext(writes)

			if err != nil {
				fmt.Println("GAK AUTH")
					next.ServeHTTP(w, r)
					return
			}

			fmt.Println("DIJEBOL")

			//validate jwt token
			tokenStr := cookie
			//fmt.Println(tokenStr.Value)
			userId, err := ParseToken(tokenStr.Value)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// create user and check if user exists in db
			user := model.User{ID: int(userId)}
			id, err := model.GetUserIdByUsername(userId)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.ID = id
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

func ForWrite(ctx context.Context) *http.ResponseWriter {
	raw, _ := ctx.Value(userCtxWrite).(*http.ResponseWriter)
	return raw
}

var (
	SecretKey = []byte("rahasia")
)

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (float64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(float64)
		return userId, nil
	} else {
		return 0, err
	}
}

