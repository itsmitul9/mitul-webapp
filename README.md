### Start the Application

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/itsmitul9/webapp.git
   cd webapp
   ```

2. **Run the Application:**
   ```bash
   chmod +x webapp
   ./webapp --port=8123
   ```

   The application will start on port 8123 by default. You can change the port by using the `--port` flag.

## API Endpoints

### `/app/status`

- **Method:** `GET`
- **Description:** Returns the current status of the application, including CPU usage and the number of replicas.
- **Response:**
  ```json
  {
      "cpu": {
          "highPriority": 0.68
      },
      "replicas": 10
  }
  ```

### `/app/replicas`

- **Method:** `PUT`
- **Description:** Updates the number of replicas based on the provided JSON payload.
- **Request Body:**
  ```json
  {
      "replicas": 11
  }
  ```
- **Response:**
$ curl -X PUT -H "Content-Type: application/json" -d '{"replicas": 11}' http://localhost:8123/app/replicas
  ```json
  {
      "message": "Replicas updated"
  }
  ```


