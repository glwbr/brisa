// Package ba implements the BA SEFAZ NFC-e portal scraper.
package ba

const (
	BaseURL         = "https://nfe.sefaz.ba.gov.br"
	AccessKeyPage   = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_chave_acesso.aspx"
	CaptchaEndpoint = "/servicos/nfce/Modulos/AntiRobo/NFCEC_anti_robo.aspx"
	DanfePage       = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_danfe.aspx"
	TabsPage        = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_abas.aspx"
)

const (
	FieldAccessKey = "txt_chave_acesso"
	FieldCaptcha   = "txt_cod_antirobo"
	FieldSubmit    = "btn_consulta_completa"
	FieldViewTabs  = "btn_visualizar_abas"
)

type Tab string

const (
	TabNFe      Tab = "nfe"
	TabProdutos Tab = "produtos"
)

func (t Tab) ButtonName() string {
	switch t {
	case TabNFe:
		return "btn_aba_nfe"
	case TabProdutos:
		return "btn_aba_produtos"
	default:
		return ""
	}
}
