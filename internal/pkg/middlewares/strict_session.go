package middlewares

import (
	"net/http"
	"os"

	"github.com/Christian-007/fit-forge-notification/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/appservices"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/security"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/utils"
)

func StrictSession(authService appservices.AuthService, secretManagerProvider security.SecretManageProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var uuid string

			cookie, err := r.Cookie("token")
			if err != nil {
				utils.SendResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Message: "Unauthorized"})
				return
			}

			privateKey, err := secretManagerProvider.GetPrivateKey(r.Context(), os.Getenv("GCP_SECRET_DIR"))
			if err != nil {
				utils.SendResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			token := cookie.Value
			claims, err := authService.ValidateToken(privateKey, token)
			if err != nil {
				if err == apperrors.ErrExpiredToken {
					utils.SendResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Message: "Token is expired"})
					return
				}

				if err == apperrors.ErrInvalidSignature {
					utils.SendResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Message: "Token is invalid"})
					return
				}

				utils.SendResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			uuid = claims.Uuid

			// It's important to use `userId` from the cache just in case the JWT has been tampered
			authData, err := authService.GetHashAuthDataFromCache(uuid)
			if err != nil {
				if err == apperrors.ErrRedisValueNotInHash {
					utils.SendResponse(w, http.StatusUnauthorized, utils.ErrorResponse{Message: "Unauthorized"})
					return
				}

				utils.SendResponse(w, http.StatusInternalServerError, utils.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			ctx := requestctx.WithUserId(r.Context(), authData.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
