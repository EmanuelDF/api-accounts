package main_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type getAccountDataStage struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

func FetchAccountTest(t *testing.T) (*getAccountDataStage, *getAccountDataStage, *getAccountDataStage) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	var stage = &getAccountDataStage{
		ID:             "",
		OrganisationID: uuid,
		Type:           "accounts",
		Version:        new(int64),
	}

	return stage, stage, stage
}

func (s *getAccountDataStage) and() *getAccountDataStage {
	return s
}

func (s *getAccountDataStage) an_authorized_service_user() *getAccountDataStage {
	s.client = NewAccountAPIClient(ServerPort)
	return s
}

func (s *getAccountDataStage) an_account_with_number_and_bank_id(accountNumber string, bankID string) *getAccountDataStage {
	s.accountNumber = accountNumber
	s.bankID = bankID

	s.createAccountsRequestModel = &models.NewAccount{
		ID:             strfmt.UUID(uuid.New().String()),
		OrganisationID: convert.FromUUID(s.organisationId),
		Type:           string(models.ResourceTypeAccounts),
		Attributes: &models.AccountAttributes{
			AccountNumber: accountNumber,
			BankID:        bankID,
			Country:       convert.StringToPtr("GB"),
		},
	}

	s.postOrganisationAccountsCreatedResult, s.error = s.client.PostOrganisationAccounts(&account_api.PostOrganisationAccountsParams{
		Context: context.Background(),
		CreationRequest: &models.AccountCreation{
			Data: s.createAccountsRequestModel,
		},
	})
	return s
}

func (s *getAccountDataStage) fetching_an_account_by_id() *getAccountDataStage {
	s.getAccountResult, s.error = s.client.GetOrganisationAccountsID(&account_api.GetOrganisationAccountsIDParams{
		Context: context.Background(),
		ID:      s.postOrganisationAccountsCreatedResult.Payload.Data.ID,
	})
	return s
}

func (s *getAccountDataStage) fetching_an_account_by_a_non_existing_id() *getAccountDataStage {
	s.getAccountResult, s.error = s.client.GetOrganisationAccountsID(&account_api.GetOrganisationAccountsIDParams{
		Context: context.Background(),
		ID:      convert.FromUUID(uuid.New()),
	})
	return s
}

func (s *getAccountDataStage) the_account_should_be_found() *getAccountDataStage {
	assert.NoError(s.t, s.error)
	account := s.getAccountResult.Payload.Data

	assert.NotNil(s.t, account)
	assert.Equal(s.t, models.ResourceTypeAccounts, account.Type)
	assert.Equal(s.t, convert.FromUUID(s.organisationId), account.OrganisationID)
	assert.Equal(s.t, s.accountNumber, account.Attributes.AccountNumber)
	assert.Equal(s.t, s.bankID, account.Attributes.BankID)
	return s
}

func (s *getAccountDataStage) the_status_code_is_404_not_found() *getAccountDataStage {
	assert.Error(s.t, s.error)
	assert.IsType(s.t, &account_api.GetOrganisationAccountsIDNotFound{}, s.error)
	return s
}
