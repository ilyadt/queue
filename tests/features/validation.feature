Feature: Validation constraints

    Scenario: Get value with negative timeout
        When I request value with negative timeout
        Then I get 400 status
