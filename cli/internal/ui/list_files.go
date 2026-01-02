package ui

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/list"
)

func ListFiles(filesData []list.FilesMetadata) {

	if len(filesData) == 0 {
		fmt.Println("No files found.")
		return
	}

	fmt.Println()
	fmt.Println("Your files:")
	fmt.Println("--------------------------------------------------------------------------------------")
	fmt.Printf("%-10s  %-25s  %-10s  %-10s  %-8s  %-12s\n",
		"ID",
		"NAME",
		"SIZE",
		"STATUS",
		"KEY",
		"CREATED",
	)
	fmt.Println("--------------------------------------------------------------------------------------")

	for _, d := range filesData {
		fmt.Printf("%-10s  %-25s  %-10s  %-10s  %-8s  %-12s\n",
			d.ID.String()[:8],
			truncate(d.FileName, 25),
			formatBytes(d.EncryptedSizeBytes),
			d.Status,
			d.KeyManagementMode,
			d.CreatedAt.Format("2006-01-02"),
		)
	}

	fmt.Println("--------------------------------------------------------------------------------------")
	fmt.Println("Use `hashdrop files show <id>` to see file details")
	fmt.Println()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div := int64(unit)
	exp := 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div),
		"KMGTPE"[exp],
	)
}
