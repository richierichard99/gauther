# Gauther

A simple Go HTTP authentication server using JWT and RSA keys.

## Features
- `/login` endpoint: Accepts username and password, returns a JWT if valid.
- `/validate` endpoint: Protected by JWT middleware, validates the provided JWT.

## Usage

### 1. (OPTIONAL) Generate an RSA Private Key
Generate a 4096-bit RSA private key in PEM format:

```sh
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096
```

### 2. (OPTIONAL) Base64-encode the Private Key and Set the Environment Variable

```sh
export GUATH_PRIVATE_KEY="$(base64 < private.pem)"
```

you can run the server without this, and it will randomly generate an rsa private key each time, this however means that after a restart previously issued jwt tokens will not be valid.

It is _heavily_ encouraged that you supply your own key if you want to use this in a production environment. Without it you also won't be able to run multiple validation servers in parellel.

### 3. Run the Server

```sh
go run ./cmd/main.go
```

The server will start on `:8080` by default.

### 4. Example Requests

#### Login (get JWT)
```sh
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

#### Validate (use JWT from login)
Replace `<TOKEN>` with the JWT you received from the login response.
```sh
curl -X GET http://localhost:8080/validate \
  -H "Authorization: <TOKEN>"
```

## Notes
- The username and password are hardcoded as `admin` and `password123`.
- The private key must be a PEM-encoded PKCS1 or PKCS8 RSA key.
- The environment variable must be base64-encoded to preserve formatting.

---
MIT License
