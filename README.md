# 💸 Amartha Loan Service API

A backend API for managing the full lifecycle of a loan from application to approval, investment, and disbursement.  
Built with **Go**, **Gin**, and Testify.

This was developed as part of the Amartha engineering coding test.

---

## 🚀 Features

- Submit new loan applications
- Approve loans with validator info
- Accept multiple investor contributions
- Disburse approved loans with agreement files
- Get individual or full loan list
- Mock email notifications for investors

---

## 🧰 Tech Stack

- Language: **Golang (1.21+)**
- Framework: **Gin**
- Testing: `testing`, `httptest`, `testify`
- Mock Email: `log.Printf()` + `fmt.Println()`
- Dependency Management: `go mod`
- Tooling: `Makefile`

---

## 📂 Project Structure

```
├── api/                # HTTP handlers and routes
├── core/loan/          # Business logic (state machine, models, service, repo)
├── email/              # MockEmailSender (logs email sends)
├── cmd/                # Main application entrypoint
├── go.mod / go.sum
├── Makefile            # Dev & CI tasks
└── README.md
```

---

## ⚙️ How to Run

```bash
go run main.go
```

---

## 🧪 How to Test
Run all tests:
```
make test
```
Check coverage in terminal:
```
make coverage
```
Open full HTML coverage report:
```
make coverage-html
```

---

## 🛠 Sample API Endpoints
```http
POST /loans
POST /loans/:id/approve
POST /loans/:id/invest
POST /loans/:id/disburse
GET  /loans/:id
GET  /loans
```

Test with Postman Collection (file in the folder)

---

## 🔖 Author
Rifqi Fauzan Akram  
Email: rifqiakram57@gmail.com  
GitHub: [@rifqiakrm](https://github.com/rifqiakrm)