package core

import (
	"testing"
)

func TestNormalizeAccessKey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "already normalized",
			input: "29250306057223031484650140003829591141073162",
			want:  "29250306057223031484650140003829591141073162",
		},
		{
			name:  "with spaces",
			input: "2925 0306 0572 2303 1484 6501 4000 3829 5911 4107 3162",
			want:  "29250306057223031484650140003829591141073162",
		},
		{
			name:  "with dashes",
			input: "2925-0306-0572-2303-1484-6501-4000-3829-5911-4107-3162",
			want:  "29250306057223031484650140003829591141073162",
		},
		{
			name:  "mixed formatting",
			input: "2925 0306-0572 2303 1484.6501.4000.3829 5911 4107 3162",
			want:  "29250306057223031484650140003829591141073162",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeAccessKey(tt.input)
			if got != tt.want {
				t.Errorf("normalizeAccessKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsValidAccessKey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid 44 digits",
			input: "29250306057223031484650140003829591141073162",
			want:  true,
		},
		{
			name:  "too short",
			input: "2925030605722303148465014000382959114107316",
			want:  false,
		},
		{
			name:  "too long",
			input: "292503060572230314846501400038295911410731620",
			want:  false,
		},
		{
			name:  "empty",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidAccessKey(tt.input)
			if got != tt.want {
				t.Errorf("isValidAccessKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBAScrapingState_String(t *testing.T) {
	tests := []struct {
		state BAScrapingState
		want  string
	}{
		{BAStateInitial, "initial"},
		{BAStateAccessKeyPage, "access_key_page"},
		{BAStateCaptchaRequired, "captcha_required"},
		{BAStateDanfePage, "danfe_page"},
		{BAStateTabsPage, "tabs_page"},
		{BAStateComplete, "complete"},
		{BAStateError, "error"},
		{BAScrapingState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.state.String()
			if got != tt.want {
				t.Errorf("BAScrapingState.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBATabButtonName(t *testing.T) {
	tests := []struct {
		tab  BATab
		want string
	}{
		{BATabTypeNFe, BATabNFe},
		{BATabTypeEmitente, BATabEmitente},
		{BATabTypeDestinatario, BATabDestinatario},
		{BATabTypeProdutos, BATabProdutos},
		{BATabTypeTotais, BATabTotais},
		{BATabTypeTransporte, BATabTransporte},
		{BATabTypeCobranca, BATabCobranca},
		{BATabTypeInfAdicionais, BATabInfAdicionais},
		{BATab("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.tab), func(t *testing.T) {
			got := BATabButtonName(tt.tab)
			if got != tt.want {
				t.Errorf("BATabButtonName(%q) = %q, want %q", tt.tab, got, tt.want)
			}
		})
	}
}
