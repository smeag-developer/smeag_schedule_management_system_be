package utils

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sikozonpc/ecom/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidatePayload(ctx *gin.Context, v interface{}) {
	// Read body once
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "failed to read body"})
		return
	}

	// Reset body so it can be read again
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Bind JSON from buffered body
	if err := ctx.ShouldBindJSON(v); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Reset again for further processing if needed
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Validate struct
	if err := Validate.Struct(v); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func HandlesPayloadValidation(ctx *gin.Context, nid string, entity any) (bson.ObjectID, error) {

	objID, objErr := IDValidation(ctx, nid)

	if objErr != nil {
		return bson.NilObjectID, objErr
	}

	if err := ParseJSON(ctx, entity); err != nil {
		return bson.NilObjectID, err
	}

	// Validate Payload
	if err := utils.Validate.Struct(entity); err != nil {
		errors := err.(validator.ValidationErrors)
		return bson.NilObjectID, errors
	}

	return objID, nil
}

func IDValidation(ctx *gin.Context, nid string) (bson.ObjectID, error) {

	vars := ctx.Param(nid)
	objID, err := bson.ObjectIDFromHex(vars)

	if err != nil {
		return bson.NilObjectID, err
	}

	return objID, nil
}
