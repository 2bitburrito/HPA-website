package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/2bitburrito/hpa-website/internal/email"
)

func (s *Server) HandleContactForm(w http.ResponseWriter, r *http.Request) {
	log.Println("Sending Email...")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<span class="error">Invalid form submission</span>`))
		return
	}
	fmt.Printf("%+v", r.Form)

	params := email.SendEmailParams{
		SendingAddress: r.PostForm.Get("form-email"),
		SenderName:     r.PostForm.Get("form-name"),
		Message:        r.PostForm.Get("form-message"),
	}

	err = email.SendContactEmail(s.Dependencies.Aws, params)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<span class="error">Sorry, something went wrong. Please try again later.</span>`))
		return
	}
	w.WriteHeader(http.StatusPermanentRedirect)
	w.Write([]byte(`<span class="success">✅ Message sent! We'll get back to you soon.</span>`))
}
