package ui

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func UploadSuccessfulMsg(fileName, fileID, s3ObjectKey string, fileSize int64) {
	msg := fmt.Sprintf(`
================= UPLOAD SUCCESSFUL 🎉 =================

Your file has been securely encrypted and uploaded.

File Name : %s
File Size : %d bytes
File ID   : %s

Download URL:
%s%s
---------------------------------------------------------
`, fileName, fileSize, fileID, config.UrlPrefix, s3ObjectKey)
	fmt.Print(msg)
}
