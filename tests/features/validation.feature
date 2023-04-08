Feature: Validation constraints

    Scenario: Get value with negative timeout
        When I request value with negative timeout
        Then I get 400 status

    Scenario: Put value without queue name
        When I put value without queue name
        Then I get 400 status
