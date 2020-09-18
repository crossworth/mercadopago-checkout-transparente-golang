package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	MercadoPagoPublicKey   = ``
	MercadoPagoAccessToken = ``
)

func main() {
	ct := CheckoutTransparente{
		MercadoPagoPublicKey:   MercadoPagoPublicKey,
		MercadoPagoAccessToken: MercadoPagoAccessToken,
	}

	http.HandleFunc("/", ct.CheckoutPage)
	http.HandleFunc("/process_payment", ct.ProcessPayment)

	log.Println("Servidor iniciado na porta 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}

type CheckoutTransparente struct {
	MercadoPagoPublicKey   string
	MercadoPagoAccessToken string
}

type CheckoutData struct {
	MercadoPagoPublicKey string
}

func (c *CheckoutTransparente) CheckoutPage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("./checkout.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_ = tpl.Execute(w, CheckoutData{
		MercadoPagoPublicKey: c.MercadoPagoPublicKey,
	})
}

func (c *CheckoutTransparente) ProcessPayment(w http.ResponseWriter, r *http.Request) {

}
