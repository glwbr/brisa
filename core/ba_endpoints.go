package core

// BA SEFAZ NFC-e portal endpoints.
const (
	// BABaseURL is the base URL for the BA SEFAZ NFC-e portal.
	BABaseURL = "https://nfe.sefaz.ba.gov.br"
	// BAAccessKeyPage = "https://nfe.sefaz.ba.gov.br/servicos/nfce/Modulos/Geral/NFCEC_consulta_chave_acesso.aspx"

	// BAAccessKeyPage is the initial page where users enter the access key.
	BAAccessKeyPage = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_chave_acesso.aspx"

	// BACaptchaEndpoint generates captcha images.
	// Requires a timestamp parameter: ?t=<milliseconds>
	BACaptchaEndpoint = "/servicos/nfce/Modulos/AntiRobo/NFCEC_anti_robo.aspx"

	// BADanfePage is the DANFE visualization page after successful captcha.
	BADanfePage = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_danfe.aspx"

	// BATabsPage is the tabbed view page with detailed invoice information.
	BATabsPage = "/servicos/nfce/Modulos/Geral/NFCEC_consulta_abas.aspx"
)

// BA form field names for ASP.NET POST requests.
const (
	// BAFieldAccessKey is the access key input field name.
	BAFieldAccessKey = "txt_chave_acesso"

	// BAFieldCaptcha is the captcha solution input field name.
	BAFieldCaptcha = "txt_cod_antirobo"

	// BAFieldSubmit is the submit button field name for the access key form.
	BAFieldSubmit = "btn_consulta_completa"

	// BAFieldViewTabs is the button to switch to tabbed view.
	BAFieldViewTabs = "btn_visualizar_abas"

	// BAFieldOriginCall is a hidden field tracking origin.
	BAFieldOriginCall = "hd_origem_chamada"

	// Additional ASP.NET hidden fields
	BAFieldLastFocus     = "__LASTFOCUS"
	BAFieldEventTarget   = "__EVENTTARGET"
	BAFieldEventArgument = "__EVENTARGUMENT"
)

// BA tab button names for switching between tabs.
const (
	BATabNFe           = "btn_aba_nfe"
	BATabEmitente      = "btn_aba_emitente"
	BATabDestinatario  = "btn_aba_destinatario"
	BATabProdutos      = "btn_aba_produtos"
	BATabTotais        = "btn_aba_totais"
	BATabTransporte    = "btn_aba_transporte"
	BATabCobranca      = "btn_aba_cobranca"
	BATabInfAdicionais = "btn_aba_infadicionais"
)

// BATab represents a tab in the BA NFC-e portal.
type BATab string

const (
	BATabTypeNFe           BATab = "nfe"
	BATabTypeEmitente      BATab = "emitente"
	BATabTypeDestinatario  BATab = "destinatario"
	BATabTypeProdutos      BATab = "produtos"
	BATabTypeTotais        BATab = "totais"
	BATabTypeTransporte    BATab = "transporte"
	BATabTypeCobranca      BATab = "cobranca"
	BATabTypeInfAdicionais BATab = "inf_adicionais"
)

// BATabButtonName returns the button field name for a given tab.
func BATabButtonName(tab BATab) string {
	switch tab {
	case BATabTypeNFe:
		return BATabNFe
	case BATabTypeEmitente:
		return BATabEmitente
	case BATabTypeDestinatario:
		return BATabDestinatario
	case BATabTypeProdutos:
		return BATabProdutos
	case BATabTypeTotais:
		return BATabTotais
	case BATabTypeTransporte:
		return BATabTransporte
	case BATabTypeCobranca:
		return BATabCobranca
	case BATabTypeInfAdicionais:
		return BATabInfAdicionais
	default:
		return ""
	}
}
