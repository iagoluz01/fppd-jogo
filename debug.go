package main

import (
	"fmt"
	"os"
	"time"
)

// RegistrarLog adiciona uma mensagem ao arquivo de log
func RegistrarLog(formato string, args ...interface{}) {
	f, err := os.OpenFile("cliente_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	msg := fmt.Sprintf(formato, args...)
	fmt.Fprintf(f, "[%s] %s\n", timestamp, msg)
}
