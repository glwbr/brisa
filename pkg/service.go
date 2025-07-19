package juice

// import (
// 	"context"
// 	"fmt"
// 	"github.com/PuerkitoBio/goquery"
// 	"io"
// 	"net/url"
// 	"strings"
// )
//
// // API domains and endpoints
// const (
// 	BaseURL             = "http://nfe.sefaz.ba.gov.br"
// 	CaptchaPath         = "/servicos/nfce/Modulos/AntiRobo/NFCEC_anti_robo.aspx"
// 	SessionInitPath     = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_chave_acesso.aspx"
// 	BasicReceiptPath    = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_danfe.aspx"
// 	DetailedReceiptPath = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_abas.aspx"
// )
//
// // error  <input type="hidden" name="__VIEWSTATEGENERATOR" id="__VIEWSTATEGENERATOR" value="909DA7CA" />
// // error  <input type="hidden" name="__EVENTVALIDATION" id="__EVENTVALIDATION" value="/wEdAAOjER9WX6L0L9yOeWdtBeFreSDB1+m2Z7vHrR/0ATQoCNR/oa2GSDmd3sGFnu/eG1JK7t/QpOQKo1OVnHHP+nsBLR0I6AuUWOCadX82rsNvrw==" />
//
// // INFO: Constants for form submission which allows us to skip a request only to retrieve these
// // we dont neeed the blank valeus but keeping it for consistency with the webservice, who knows.
//
// // Form field names
// const (
// 	FormCaptcha               = "txt_cod_antirobo"
// 	FormAccessKey             = "txt_chave_acesso"
// 	FormViewState             = "__VIEWSTATE"
// 	FormEventTarget           = "__EVENTTARGET"
// 	FormSubmitButton          = "btn_consulta_completa"
// 	FormEventArgument         = "__EVENTARGUMENT"
// 	FormEventValidation       = "__EVENTVALIDATION"
// 	FormViewStateGenerator    = "__VIEWSTATEGENERATOR"
// 	FormHomeSubmitButtonValue = "Consultar"
// )
//
// const (
// 	TestAccessKey1 = "29250306057223031484650140003829591141073162"
// 	TestAccessKey2 = "29250206057223031484650180002493191180299387"
// 	TestAccessKey3 = "29250409182947000216651100001675181110128180"
// )
//
// // HTTP client defaults
// const (
// 	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
// 	DefaultTimeout   = 10 // seconds
// )
//
// type FormState struct {
// 	ViewState          string
// 	EventValidation    string
// 	ViewStateGenerator string
// 	HTML               RawInvoiceData
// }
//
// // GetHomePage - Fetches the initial page and saves cookies
// // maybe we can replace the Location and go directly to abas instead of:
// // servicos/nfce/modulos/geral/NFCEC_consulta_danfe.aspx
// func (c *Client) GetHomePage(ctx context.Context) (*FormState, error) {
// 	res, err := c.Get(ctx, SessionInitPath, &RequestOptions{
// 		QueryParams: url.Values{"param1": {"value1"}},
// 	})
//
// 	if err != nil {
// 		return nil, fmt.Errorf("[juicy]: error making request: %w", err)
// 	}
// 	defer res.Body.Close()
//
// 	bodyBytes, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("[juicy]: error reading body: %w", err)
// 	}
//
// 	body := string(bodyBytes)
//
// 	// TODO: move this elsewhere
// 	// if c.options.DebugMode {
// 	// 	os.WriteFile("sefaz_page.html", []byte(body), 0644)
// 	// }
//
// 	state, err := parseFormState(body)
// 	if err != nil {
// 		return nil, fmt.Errorf("[juicy]: error parsing form state: %w", err)
// 	}
// 	state.HTML = RawInvoiceData(body) // optional
//
// 	return state, nil
// }
//
// func parseFormState(html string) (*FormState, error) {
// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	getField := func(name string) string {
// 		val, _ := doc.Find(fmt.Sprintf("input[name='%s']", name)).Attr("value")
// 		return val
// 	}
//
// 	return &FormState{
// 		ViewState:          getField("__VIEWSTATE"),
// 		EventValidation:    getField("__EVENTVALIDATION"),
// 		ViewStateGenerator: getField("__VIEWSTATEGENERATOR"),
// 	}, nil
// }
//
// func (c *Client) SubmitForm(accessKey string, captchaAnswer string) (string, error) {
// 	formData := url.Values{
// 		"__EVENTARGUMENT":       {""},
// 		"__EVENTTARGET":         {""},
// 		"__EVENTVALIDATION":     {""},
// 		"__VIEWSTATE":           {""},
// 		"__VIEWSTATEGENERATOR":  {""},
// 		"txt_chave_acesso":      {""},
// 		"txt_cod_antirobo":      {""},
// 		"btn_consulta_completa": {""},
// 	}
//
// 	// req, err := http.NewRequest("POST", SEFAZEndpoint, strings.NewReader(formData.Encode()))
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("[juicy] creating home request: %v", err)
// 	// }
//
// 	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	// req.Header.Add("User-Agent", DefaultUserAgent)
//
// 	// for _, cookie := range c.cookies {
// 	// 	req.AddCookie(cookie)
// 	// }
//
// 	// res, err := c.client.Do(req)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("[juicy] error submitting form: %v", err)
// 	// }
//
// 	// defer res.Body.Close()
//
// 	// body, err := io.ReadAll(res.Body)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("[juicy] error reading response: %v", err)
// 	// }
//
// 	return string(formData["__EVENTVALIDATION"][0]), nil
// }
