package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	MercadoPagoPublicKey = ``
	MercadoPagoSecretKey = ``
)

func main() {
	ct := CheckoutTransparente{}

	http.HandleFunc("/", ct.CheckoutPage)
	http.HandleFunc("/process_payment", ct.ProcessPayment)

	log.Println("Servidor iniciado na porta 8080")
	log.Println(http.ListenAndServe(":8080", nil))
}

type CheckoutTransparente struct {
	MercadoPagoPublicKey string
}

type CheckoutData struct {
	MercadoPagoPublicKey string
}

func (c *CheckoutTransparente) CheckoutPage(w http.ResponseWriter, r *http.Request)  {
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
