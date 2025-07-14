# URL Shortener

A microservice-based URL shortener with logging middleware, built with Go (backend) and React (frontend).

## Project Structure

```
12220607/
  frontend/   # React + Vite + TS (UI)
  backend/    # Go microservice (API)
  logging/    # Go logging middleware
```

## Running the App

1. **Install dependencies:**
   - Frontend: `cd frontend && npm install`
   - Backend: `cd backend && go mod tidy`

2. **Start the backend (port 3001):**
   ```sh
   cd backend
   go run main.go
   ```

3. **Start the frontend (port 3000):**
   ```sh
   cd frontend
   npm run dev
   ```
   The app will be available at [http://localhost:3000](http://localhost:3000)

---

## API Endpoints

### Create Short URL
- **POST** `/shorturls`
- **Request Body:**
  ```json
  {
    "url": "https://leetcode.com/contest/biweekly-contest-161/",
    "validity": 60,           // (optional, seconds, default 30)
    "shortcode": "custom1"   // (optional)
  }
  ```
- **Response:**
  ```json
  {
    "shortcode": "abc123",
    "expiry": "2025-07-14T12:34:56+05:30"
  }
  ```

### Get Short URL Stats
- **GET** `/shorturls/{shortcode}`
- **Response:**
  ```json
  {
    "url": "https://leetcode.com/contest/biweekly-contest-161/",
    "created_at": "2025-07-14T12:34:56+05:30",
    "expiry": "2025-07-14T12:35:26+05:30",
    "hits": 2,
    "clicks": [
      {
        "timestamp": "2025-07-14T12:35:10+05:30",
        "referrer": "",
        "location": "127.0.x.x"
      }
    ]
  }
  ```

### Redirect
- **GET** `/{shortcode}`
- Redirects to the original URL if valid and not expired.

---

## App Screenshots

> ![App Home](images/app-home.png)
> ![Stats Example](images/app-stats.png)

(Replace with your own screenshots after running the app)

---

## Logging
- All backend requests and errors are logged via the logging middleware in the `logging/` folder.

---

## Deployment
- Push only the `frontend/`, `backend/`, and `logging/` folders to GitHub.
- `.gitignore` is included to keep your repo clean. 