Feature: Validation constraints

    Scenario: Get value with negative timeout
        When I request value with negative timeout
        Then I get 400 status

    Scenario: Put value without queue name
        When I put value "some_value" in queue ""
        Then I get 400 status

    Scenario: Put empty value in queue
        When I put value "" in queue "some_queue"
        Then I get 400 status
