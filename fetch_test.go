package main_test

import "testing"

func TestAcc_FetchAccount(t *testing.T) {
	given, when, then := FetchAccountTest(t)

	given.
		an_authorized_service_user().and().
		an_account_with_number_and_bank_id("41426819", "400300")

	when.
		fetching_an_account_by_id()

	then.
		the_account_should_be_found()
}

func TestAcc_FetchAccount_NotFound(t *testing.T) {
	given, when, then := FetchAccountTest(t)

	given.
		an_authorized_service_user()

	when.
		fetching_an_account_by_a_non_existing_id()

	then.
		the_status_code_is_404_not_found()
}
