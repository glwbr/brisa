package main

import (
	"context"
	"fmt"
	"io"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/glwbr/brisa/internal/client"
	"github.com/glwbr/brisa/pkg/errors"
	"github.com/glwbr/brisa/pkg/logger"
)

var (
	captchaPath     = "/servicos/nfce/Modulos/AntiRobo/NFCEC_anti_robo.aspx"
	baseURL         = "http://nfe.sefaz.ba.gov.br/servicos/nfce/"
	sessionInitPath = "Modulos/Geral/NFCEC_consulta_chave_acesso.aspx"
)

func main() {
	l := logger.NewStandardLogger(os.Stdout, logger.DebugLevel)

	// handle cases where baseURL have a trainlling slash
	// where it doesnt have
	// where paths have a prefix slash and when it doesnt have
	// if the baseURL is provided we should cleanup // nroamlize it to have or not have the trailling slash?
	// then when the path is provided we can expect a default behaviour
	jar, _ := cookiejar.New(nil)
	c, err := client.New(client.WithLogger(l), client.WithDebug(true), client.WithBaseURL(baseURL), client.WithCookieJar(jar))
	if err != nil {
		l.Error("Failed to create client", "error", err)
		os.Exit(1)
	}

	res, err := c.Get(context.Background(), sessionInitPath, nil)
	if err != nil {
		l.Error("API request failed", "error", err)

		if errors.IsNotFound(err) {
			l.Info("page not found", err)
		}

		os.Exit(1)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	// parse or extract ?
	s, _ := parseFormState(string(body))
	pp := logger.PrettyLogger{Logger: l}
	pp.InfoPretty("Form state extracted", *s)

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// generate the captcha image
	// the sessionCookie must be attached this is being provided by the cookiejar
	captchRes, _ := c.Get(context.Background(), captchaPath, &client.RequestConfig{
		Params: url.Values{
			"t": {strconv.FormatInt(timestamp, 10)},
		},
	})

	defer captchRes.Body.Close()

	imgBytes, _ := io.ReadAll(captchRes.Body)
	name := "captcha_" + strconv.FormatInt(timestamp, 10) + ".jpeg"
	err = os.WriteFile(name, imgBytes, 0o644)
	if err != nil {
		// TODO: use a custom error CaptchaError
		l.Error("error while generating the captcha", err)
	}

	var captchaAnswer string
	_, err = fmt.Scanln(&captchaAnswer)
	if err != nil {
		l.Error("error reading the answer", err)
	}

	submitForm(TestAccessKey1, captchaAnswer, c, s, timestamp)
}

const (
	TestAccessKey1 = "29250306057223031484650140003829591141073162"
	TestAccessKey2 = "29250206057223031484650180002493191180299387"
	TestAccessKey3 = "29250409182947000216651100001675181110128180"
)

// these are JUST A TEST WAY OF DOING TO CHECK ITS ACTUALLY WORKING
// IT WORKS !!
func submitForm(id, captchaAnswer string, c *client.Client, s *FormState, timestamp int64) (string, error) {
	formData := url.Values{
		"__EVENTARGUMENT":       {""},
		"__EVENTTARGET":         {""},
		"__EVENTVALIDATION":     {s.EventValidation},
		"__VIEWSTATE":           {s.ViewState},
		"__VIEWSTATEGENERATOR":  {s.ViewStateGenerator},
		"txt_chave_acesso":      {id},
		"txt_cod_antirobo":      {captchaAnswer},
		"btn_consulta_completa": {"Consultar"},
	}

	// we do a post to the same initial endpoint, now with the form values set
	// maybe we should also rename the var to better reflect these, maybe initPath ?
	// this redirects us to the basic receipt enpoint
	// BasicReceiptPath    = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_danfe.aspx"
	// if everything is okay here we already have some information for the invoice
	// but maybe we want to cancel the redirect and go directly to the detailed view
	// DetailedReceiptPath = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_abas.aspx"
	res, err := c.Post(context.Background(), sessionInitPath, &client.RequestConfig{
		Body:    strings.NewReader(formData.Encode()),
		Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	})
	if err != nil {
		return "", fmt.Errorf("[brisa] error submitting form: %v", err)
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("[brisa] error reading response: %v", err)
	}

	name := "raw_invoice" + strconv.FormatInt(timestamp, 10) + ".html"
	os.WriteFile(name, bodyBytes, 0o644)

	return string(bodyBytes), nil
}

type FormState struct {
	ViewState          string
	EventValidation    string
	ViewStateGenerator string
}

func parseFormState(html string) (*FormState, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	getField := func(name string) string {
		val, _ := doc.Find(fmt.Sprintf("input[name='%s']", name)).Attr("value")
		return val
	}

	return &FormState{
		ViewState:          getField("__VIEWSTATE"),
		EventValidation:    getField("__EVENTVALIDATION"),
		ViewStateGenerator: getField("__VIEWSTATEGENERATOR"),
	}, nil
}
