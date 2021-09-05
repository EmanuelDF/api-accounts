package models

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.

type NewAccountData struct {
	Type           string                `json:"type,omitempty"`
	ID             string                `json:"id,omitempty"`
	OrganisationID string                `json:"organisation_id,omitempty"`
	Attributes     *NewAccountAttributes `json:"attributes,omitempty"`
}

type NewAccountAttributes struct {
	Country                *string `json:"country,omitempty"`
	BaseCurrency           string  `json:"base_currency,omitempty"`
	BankID                 string  `json:"bank_id,omitempty"`
	BankIDCode             string  `json:"bank_id_code,omitempty"`
	Bic                    string  `json:"bic,omitempty"`
	ProcessingService      string  `json:"processing_service,omitempty"`
	UserDefinedInformation string  `json:"user_defined_information,omitempty"`
	ValidationType         string  `json:"validation_type,omitempty"`
	ReferenceMask          string  `json:"reference_mask,omitempty"`
	AcceptanceQualifier    string  `json:"acceptance_qualifier,omitempty"`
}
