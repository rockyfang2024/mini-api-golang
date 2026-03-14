# Deployment Instructions

1. **Clone the Repository**
   ```bash
   git clone https://github.com/rockyfang2024/mini-api-golang.git
   cd mini-api-golang
   ```

2. **Install Dependencies**
   Make sure you have Go installed. Run:
   ```bash
   go mod tidy
   ```

3. **Build the Application**
   ```bash
   go build -o mini-api
   ```

4. **Run the Application**
   ```bash
   ./mini-api
   ```

5. **Access the API**
   Open your browser or use a tool like Postman to access the API endpoints at `http://localhost:8080`.

# Access Instructions

You can access the API endpoints using any REST client. Here are a few example endpoints:
- `GET /api/v1/resources` - Fetch all resources
- `POST /api/v1/resources` - Create a new resource

Remember to review the API documentation for authentication and other requirements.

---

**Current Date and Time (UTC): 2026-03-14 04:13:58**  
**Current User's Login: rockyfang2024**