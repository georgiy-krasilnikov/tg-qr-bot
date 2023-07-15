package qr

import (
	"fmt"

	"github.com/google/uuid"
	qr "github.com/skip2/go-qrcode"
)

func CreateQR(url string) (string, error) {
	fileName := uuid.New().String() + ".png"
	if err := qr.WriteFile(url, qr.Medium, 256, fileName); err != nil {
		return "", fmt.Errorf("failed to create QR file: %s", err.Error())
	}

	return fileName, nil
}
