Feature: Queue FIFO

  Scenario Template: Put N elements in queue, read them consequently
    Given there are <n> elements in queue in order from one to N
    When I get <get> elements from queue
    Then next element will be <next>

    Examples:
      | n     | get   | next |
      | 1000  |   5   | 6    |
      | 1000  |   0   | 1    |
      | 1000  |   999 | 1000 |


  Scenario: Multiple subscribers for one value must be ordered fifo
    Given 100 subscribers waiting for value in queue
    When I put enough elements in queue
    Then subscribers get values in the fifo order


  Scenario: Cancel requests. Sometimes client cancels requests so it must not be dropped from queue
    Given 10 subscribers waiting for value in queue
    When 3 subscribers cancel request
    When 7 elements pushed to queue
    Then 7 subscribers got values
    Then Queue is empty
