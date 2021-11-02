Feature: Account creation

  Scenario: Create new account
    Given a new account is set up from file "tests/resources/account.json"
    When the user uses the client to create it
    Then the user should get a response with the same account plus extra fields

  Scenario: Create new account with missing fields
    Given a new account is set up from file "tests/resources/account-missing-fields.json"
    When the user uses the client to create it
    Then the user should get the following error:
      """
      error: validation failure list:
      country in body should match '^[A-Z]{2}$'
      organisation_id in body is required
      type in body is required
      """
      
  Scenario: Create new account with empty fields
    Given a new account is set up from file "tests/resources/account-empty-fields.json"
    When the user uses the client to create it
    Then the user should get the following error:
      """
      error: validation failure list:
      account_classification in body should be one of [Personal Business]
      country in body should match '^[A-Z]{2}$'
      id in body is required
      organisation_id in body is required
      type in body is required
      """