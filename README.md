# points
Points backend app written in Go.

# Running and using the app:
Requirements: go, make, and Postman or another endpoint testing software

Run the following command to start the backend:
```bash
make run
```

# Running the app in Docker:

From the root directory, run the following command:
```bash
docker-compose up
```

# Calling the API:

Endpoints:
```
-http://localhost:8081/api/add-transaction
-http://localhost:8081/api/spend-points
-http://localhost:8081/api/get-balance
```

Example Calls:
```
Add Transaction:
<img width="1077" alt="Screenshot 2022-10-31 at 08 59 33" src="https://user-images.githubusercontent.com/44103767/199025660-1a51322b-21e2-4e8d-aedc-b3800840f7c7.png">

Spend Points:
<img width="1081" alt="Screenshot 2022-10-31 at 09 00 19" src="https://user-images.githubusercontent.com/44103767/199025856-ec7fd6d0-63ab-4631-a303-7f05ede60ca6.png">

Get Balance:
<img width="1082" alt="Screenshot 2022-10-31 at 09 01 06" src="https://user-images.githubusercontent.com/44103767/199026052-ab2e5c29-cf21-48af-859f-8873a8c5bfb3.png">
