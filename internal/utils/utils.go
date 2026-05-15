package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	common "nxt_match_event_manager_api/internal/models/common"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, common.ErrorResponse{
		ErrorCode: status,
		TimeStamp: time.Now(),
		ErrorMsg:  err.Error(),
	})
}

func ParseJSON(r *gin.Context, v interface{}) error {
	if r.Request.Body == nil {
		return fmt.Errorf("invalid request body")
	}

	return json.NewDecoder(r.Request.Body).Decode(v)
}

func GetTokenFromRequest(w http.ResponseWriter, r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	const prefix = "Bearer "

	if tokenAuth == "" {
		WriteError(w, http.StatusForbidden, fmt.Errorf("missing authorization header"))
		return ""
	}

	if !strings.HasPrefix(tokenAuth, prefix) {
		WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token format"))
		return ""
	}

	token := strings.TrimPrefix(tokenAuth, prefix)

	return token
}

func LoggerHandler(level slog.Level, message string, method string, url_path string, status int) {

	// Define a logger with a custom format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level.Level(),
	}))
	slog.SetDefault(logger)

	slog.LogAttrs(
		context.Background(),
		level,
		message,
		slog.Group("request",
			slog.String("method ", method),
			slog.String("url", url_path)),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
		slog.Int("HttpStatus", status),
	)
}

func StructToBsonM(data any) (bson.M, error) {
	bsonBytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}

	var rawMap bson.M
	if err := bson.Unmarshal(bsonBytes, &rawMap); err != nil {
		return nil, err
	}

	for k, v := range rawMap {
		if isZero(v) {
			delete(rawMap, k)
		}
	}
	return rawMap, nil
}

func isZero(v any) bool {
	switch val := v.(type) {
	case string:
		return val == ""
	case int, int32, int64, float64:
		return val == 0
	case bool:
		return !val
	case bson.ObjectID:
		return val == bson.NilObjectID
	case time.Time:
		return val.IsZero()
	case bson.DateTime:
		// Convert to time.Time and check zero
		return time.UnixMilli(int64(val)).IsZero()
	case *time.Time:
		return val == nil || val.IsZero()
	case nil:
		return true
	default:
		return false
	}
}

func MustHex(id string) (bson.ObjectID, error) {
	if id == bson.NilObjectID.Hex() {
		return bson.NilObjectID, fmt.Errorf("invalid hex id")
	}
	return bson.ObjectIDFromHex(id)
}

func IsValidSHA256(s string) error {
	if len(s) != 64 {
		return fmt.Errorf("invalid SHA256 id")
	}
	_, err := hex.DecodeString(s)

	if err != nil {
		return fmt.Errorf("unable to decode string :%v", err)
	}

	return nil
}

func MustSHA256(s string) (string, error) {

	h := sha256.New()
	h.Write([]byte(s))
	hashed := hex.EncodeToString(h.Sum(nil))

	if len(hashed) != 64 {
		return "", fmt.Errorf("invalid SHA256 id")
	}

	return hashed, nil
}

func MustStringSHA256(s string) string {

	h := sha256.New()
	h.Write([]byte(s))
	hashed := hex.EncodeToString(h.Sum(nil))

	return hashed
}

func ValidateMultipleHexId(ids []string) ([]bson.ObjectID, error) {

	hexConvert := make([]bson.ObjectID, len(ids))

	// Validate hex ids
	for i, id := range ids {
		hex, err := MustHex(id)
		if err != nil {
			slog.Error("err conversation hex", "index", i, "error", err)
			return []bson.ObjectID{}, err
		}
		hexConvert[i] = hex
	}

	return hexConvert, nil
}

func BoolPointer(b bool) *bool {
	return &b
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Marshal to JSON
	data, err := json.Marshal(obj)
	if err != nil {
		slog.Error("unable to convert struct to map", "error", err)
		return nil, err
	}

	// Unmarshal back to map
	err = json.Unmarshal(data, &result)
	return result, err
}

func SliceInterfaceToStruct[T any](data []interface{}) ([]T, error) {
	result := make([]T, 0, len(data))

	for i, item := range data {
		if item == nil {
			continue
		}

		var raw []byte
		switch v := item.(type) {
		case []byte:
			raw = v
		case string:
			raw = []byte(v)
		default:
			return nil, fmt.Errorf("element at index %d has unexpected type %T", i, item)
		}

		// Try array first (your case)
		var arr []T
		if err := json.Unmarshal(raw, &arr); err == nil {
			result = append(result, arr...)
			continue
		}

		// Fallback: try single object
		var single T
		if err := json.Unmarshal(raw, &single); err != nil {
			return nil, fmt.Errorf("failed to unmarshal element at index %d: %w", i, err)
		}
		result = append(result, single)
	}

	return result, nil
}
