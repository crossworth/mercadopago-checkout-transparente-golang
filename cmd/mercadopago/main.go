package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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

type ErrorData struct {
	Error string
}

func sendError(w http.ResponseWriter, errorMessage string) {
	tpl, err := template.ParseFiles("./error.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_ = tpl.Execute(w, ErrorData{
		Error: errorMessage,
	})
}

var (
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func (c *CheckoutTransparente) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	transactionAmount := r.FormValue("transactionAmount")
	description := r.FormValue("description")
	installments := r.FormValue("installments")
	paymentMethodID := r.FormValue("paymentMethodId")
	issuerID := r.FormValue("issuer")
	payerEmail := r.FormValue("email")

	if token == "" {
		sendError(w, "Token não informado")
		return
	}

	if transactionAmount == "" {
		sendError(w, "Valor não informado")
		return
	}

	if installments == "" {
		sendError(w, "Número de parcelas não informado")
		return
	}

	if paymentMethodID == "" {
		sendError(w, "Método de pagamento não informado")
		return
	}

	if issuerID == "" {
		sendError(w, "Emissor não informado")
		return
	}

	if payerEmail == "" {
		sendError(w, "Email do comprador não informado")
		return
	}

	type requestData struct {
		TransactionAmount int    `json:"transaction_amount"`
		Token             string `json:"token"`
		Description       string `json:"description"`
		Installments      int    `json:"installments"`
		PaymentMethodID   string `json:"payment_method_id"`
		IssuerID          int    `json:"issuer_id"`
		Payer             struct {
			Email string `json:"email"`
		} `json:"payer"`
	}

	data, err := json.Marshal(requestData{
		TransactionAmount: ToInt(transactionAmount),
		Token:             token,
		Description:       description,
		Installments:      ToInt(installments),
		PaymentMethodID:   paymentMethodID,
		IssuerID:          ToInt(issuerID),
		Payer: struct {
			Email string `json:"email"`
		}{
			Email: payerEmail,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.mercadopago.com/v1/payments?access_token=%s", MercadoPagoAccessToken), bytes.NewReader(data))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Context-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer res.Body.Close()


	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write(out)
}

func ToInt(i string) int {
	r, _ := strconv.Atoi(i)
	return r
}
