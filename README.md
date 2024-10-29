# assignment

We want to do a campaign for Uniswap's users. In the campaign, we will have 2 tasks for users to collect points. At the end of the campaign, users will get rewards based on how many points they accumulate.
This campaign will run for 4 weeks. Tasks are calculated on a weekly basis, except for the onboarding task.
The start time of this campaign can be set by the database or a config file. This allows us to set the campaign status to either a past backtest mode or a current active mode that can receive real-time events.

### Tasks
1. Onboarding task
    - The user needs to swap at least 1000u
    - Users will get 100 points immediately when they complete it
2. Share Pool Task
    - Based on the propotion of user's swap volume (USD) among all users on pool
    - Users coluld share the points pool of 10,000 points at the end of the task
    - The user needs to complete the onboarding task when we distribute the points
    - No need to be real-time

### APIs
1. Get user tasks status (is completed, amount, points, etc) by address, campaign_id
2. Get user points history for distrubuted tasks
3. Get a leaderboard based on points of distributed tasks

---

## Getting Started

1. Prepare development environment
    ```bash
    make uninstall-bin install-bin mock
    ```
2. Deploy infra using Docker Compose
    ```bash
    cd setup/infra
    docker compose up -d
    ```
3. Run server `go run cmd/server.go`
4. Try to trigger APIs as below:
    - Do swap
        ```bash
        curl --location 'http://localhost:8888/v1/swap' \
        --header 'Content-Type: application/json' \
        --data '{
            "userID": "csf4eh04v3c4kmbslpt0",
            "pair": "USDC/ETH",
            "amount": 1000
        }'
        ```
    - Get user tasks status
        ```bash
        curl --location 'http://localhost:8888/v1/addresses/0x1234567890abcdef1234567890abcdef12345678/campaigns/esf4eh04v3c4kmbslpt0'
        ```
    - Get user points history for distrubuted tasks
        ```bash
        curl --location 'http://localhost:8888/v1/users/csf4eh04v3c4kmbslpt0/pointHistories?size=10'
        ```
    - Get a leaderboard based on points of distributed tasks
        ```bash
        curl --location 'http://localhost:8888/v1/campaigns/esf4eh04v3c4kmbslpt0/leaderboard'
        ```
