
package main 

import (
	"log/slog"
	"os"
)

func main(){

	logger:=slog.New(slog.NewJSONHandler(os.Stdout,nil))

	logger.Info("This is a structured log message",

	           "copypastehero",true ,"breakmaster","vit batchelor")

	slog.Info("Login attempt","ip", 2000 ,"jamesbond",5000)
}