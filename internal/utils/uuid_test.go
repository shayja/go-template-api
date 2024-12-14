package utils_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shayja/go-template-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsValidUUID(t *testing.T) {
	// Test a valid UUID
	validUUID := uuid.New().String()
	assert.True(t, utils.IsValidUUID(validUUID), "Expected valid UUID to return true")

	// Test an invalid UUID
	invalidUUID := "invalid-uuid-string"
	assert.False(t, utils.IsValidUUID(invalidUUID), "Expected invalid UUID to return false")

	// Test an empty string
	emptyUUID := ""
	assert.False(t, utils.IsValidUUID(emptyUUID), "Expected empty string to return false")
}

func TestCreateNewUUID(t *testing.T) {
	// Generate a new UUID
	newUUID := utils.CreateNewUUID()

	// Validate that it's a proper UUID
	assert.NotEmpty(t, newUUID, "Expected generated UUID not to be empty")
	assert.True(t, utils.IsValidUUID(newUUID.String()), "Expected generated UUID to be valid")

	// Check uniqueness by generating multiple UUIDs
	anotherUUID := utils.CreateNewUUID()
	assert.NotEqual(t, newUUID, anotherUUID, "Expected generated UUIDs to be unique")
}
