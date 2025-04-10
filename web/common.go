package web

import (
	"fmt"
	"html/template"
	"net/http"
)

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
