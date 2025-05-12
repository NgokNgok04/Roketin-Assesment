# 🎬 Roketin Assessment - Converter Time and Movie API

This repository contains two challenges built using **Go**, demonstrating the use of HTTP servers, file uploads, data persistence with PostgreSQL, and API documentation using Swagger.

---

## 🚀 Features

### ✅ Challenge 1 – Time Converter

- Convert Earth time to Roketin Planet time using a custom logic handler.

### ✅ Challenge 2 – Movie Management API

- Create and upload a movie with a video file.
- Update movie details and optionally replace the video.
- Search movies by title, description, artist name, or genre.
- Paginated movie list with `page` and `limit` query params (`default: page=1, limit=10`).
- Automatically create artists and genres if they don't already exist.
- Swagger UI for interactive API documentation.

---

## 🏃‍♂️ HOW TO RUN

1. Clone the Repository

```
git clone https://github.com/your-username/roketin-assessment.git
cd roketin-assessment
```

2. Copy .env.example and create .env file in the root folder
3. Install Go Dependencies

```
go mod tidy
```

4. Run **Challenge 1 (Time Converter)**

```
go run .\Challenge_1\main.go
```

5. For **Challenge 2 (Movie API)**

✅ Start PostgreSQL via Docker

```
docker-compose up --build
```

✅(Optional) If you want dummy data. Seed dummy data

```
go run .\Challenge_2\seed\seed.go
```

✅ Run API Server in **Root Folder**

```
go run .\Challenge_2\main.go
```

The server will run at:

```
http://localhost:3000
```

6. Access Swagger API Documentation

```
http://localhost:3000/swagger/index.html/
```

If you see an error, ensure the Swagger docs are generated. Run:

```
cd .\Challenge_2\
swag init -g main.go -o docs
cd..
go run .\Challenge_2\main.go
```
