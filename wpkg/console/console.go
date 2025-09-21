package console

import "fmt"

func Success(msg string) {
	fmt.Printf("✅ %s\n", msg)
}

func Error(msg string) {
	fmt.Printf("❌ %s\n", msg)
}

func Info(msg string) {
	fmt.Printf("ℹ️  %s\n", msg)
}

func Warn(msg string) {
	fmt.Printf("⚠️  %s\n", msg)
}
