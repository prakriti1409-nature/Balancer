# Go Load Balancer

A simple **HTTP load balancer written in Go** with round-robin routing, logging, and a minimal web UI.

---
## Deployed Site :
-https://balancer-2.onrender.com/

---
## Features

- **Round-robin load balancing** between multiple backend servers
- **Health Check UI**  
  - `/health` – Displays a styled “Healthy” page
- **Metrics Dashboard**  
  - `/stats` – Shows requests handled by each backend
- **Request Logging**  
  - All requests are logged to `requests.log` (method, path, backend)
- **Web UI**  
  - Homepage with buttons to view Health and Stats

---

> **Note:** `backend.go` is only for local testing and not needed for deployment.

---

## How It Works

1. Requests come to the load balancer (`/` or any other route).
2. They are distributed to backend servers using a **round-robin algorithm**.
3. The balancer:
   - Logs each request
   - Updates backend request counts
   - Serves `/health` and `/stats` without forwarding

---

## Run Locally

### 1. Start dummy backends

Open 3 terminals and run:

```bash
go run backend.go 8081
go run backend.go 8082
go run backend.go 8083
```

```bash
go run loadbalancer.go
```
## Endpoints

- **[`/`](http://localhost:8080/)**  
  Homepage with buttons

- **[`/health`](http://localhost:8080/health)**  
  Health status (UI)

- **[`/stats`](http://localhost:8080/stats)**  
  Requests count per backend

- **Any other path** (e.g., `/api`, `/test`)  
  Will be proxied to a backend


<img width="1916" height="374" alt="image" src="https://github.com/user-attachments/assets/08a7ae0e-a7e5-4a3d-b1f1-5927ec11cd2a" />
<img width="1916" height="469" alt="image" src="https://github.com/user-attachments/assets/6ec49975-5120-4fb1-8e3e-26b1ff323cdd" />
<img width="1918" height="441" alt="image" src="https://github.com/user-attachments/assets/604ce563-3d92-4f59-9315-c7b76bdb0d63" />

