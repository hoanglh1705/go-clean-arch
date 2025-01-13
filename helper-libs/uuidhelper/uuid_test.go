package uuidhelper

import (
	"fmt"
	"testing"
)

func TestNewUuidV7Success(t *testing.T) {
	uuidV7Value := NewUuidV7String()
	fmt.Println("uuid v7: ", uuidV7Value)
	if uuidV7Value == "" {
		t.Error("failed to gen uuid v7")
	}
}
