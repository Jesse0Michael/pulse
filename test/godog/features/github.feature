Feature: github cli command
  Pulse users can use the github cli command generate AI empowered insights on github activity

  Background:
    Given I have a clean rest assured environment

  Scenario: github summary for user
    Given rest assured returns "testdata/assured/github/listevents.json" on a "GET" to "/users/jesse0michael/events"
    And rest assured returns "testdata/assured/openai/chatcompletion.json" on a "GET" to "/chat/completions"
    When I run the pulse github command on user "jesse0michael"
    Then the pulse output should equal "testdata/github/expected.txt"
