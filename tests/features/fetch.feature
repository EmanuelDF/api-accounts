Feature: Account deletion

  Scenario: Create new account, then fetch it and then delete it
    Given a new account is set up from file "tests/resources/account.json"
    When the user uses the client to create it
    And the user uses the client to fetch the created account
    Then the fetched data should exactly match the created data
    And the user uses the client to delete the created account

  Scenario: Fetch an incorrect account id
    Given an account ID of "incorrect-account-id"
    When the user uses the client to fetch it
    Then the user should get the error:
      """
      error: id is not a valid uuid
      """