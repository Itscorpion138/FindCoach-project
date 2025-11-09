Here are all the steps, sequentially, from zero to running your Go app locally:

---

### 1. Install prerequisites

* **Docker & Docker Compose**

  ```bash
  sudo apt update
  sudo apt install docker.io docker-compose -y
  sudo systemctl enable --now docker
  ```
* **Go** (if not installed)

  ```bash
  wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
  source ~/.profile
  ```

---

### 2. Clone project (if applicable)

```bash
git clone <your-repo-url>
cd gymProjPractice
```

---

### 3. Create environment file
NO NEED

### 4. Docker Compose for PostgreSQL

NO NEED

### 5. Start PostgreSQL container

```bash
docker compose up -d
```

* Check container status:

```bash
docker ps
```

---

### 6. Wait for DB to initialize

* Ensure Postgres is running and listening on port 5432.
* You can test:

```bash
docker exec -it <container_name> psql -U youruser -d yourdb
```

---

### 7. Install Go dependencies

```bash
go mod tidy
```

---

### 8. Run the Go application

```bash
go run main.go
```

* Follow the CLI prompts:

  * `addUser`
  * `deleteUser`
  * `editUser`

