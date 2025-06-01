package web

import (
	"log"
	"fmt"
	"html/template"
	"net/http"

	"github.com/stracker-go/pkg"
)

func CheckPassword(r *http.Request) bool {
	password := r.Header.Get("password")
	env, err := pkg.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	if password != env.Password {
		return false
	}
	return true
}

func WriteConfetti(w http.ResponseWriter, emoji string) error {
  t, err := template.ParseFiles("web/templates/confetti.html")
  if err != nil {
    return fmt.Errorf("Failed to parse confetti template: ", err)
  }
  err = t.Execute(w, emoji)
  if err != nil {
    return fmt.Errorf("Failed to execute confetti template: ", err)
  }
  return nil
}

func WriteCloseModal(w http.ResponseWriter) error {
  t, err := template.ParseFiles("web/templates/close-modal.html")
  if err != nil {
    return fmt.Errorf("Failed to parse close-modal template: ", err)
  }
  err = t.Execute(w, nil)
  if err != nil {
    return fmt.Errorf("Failed to execute close-modal template: ", err)
  }
  return nil
}
