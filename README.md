# NonAuth - Covert Mutual Authentication TLS Extension

**NonAuth** is a production-ready implementation of the protocol described in [draft-nonauth-00](#) (“The NonAuth Protocol: An Adaptive Anti-Censorship TLS Extension for Covert Authentication”).  
It enables covert, authenticated TLS handshakes by embedding cryptographic authentication material inside the `ClientHello` and `ServerHello` random fields, making mutual authentication indistinguishable from ordinary TLS to external observers.

---

## Features

- **Covert Mutual Authentication**  
  Embeds authentication data within TLS handshake random fields using ChaCha20-Poly1305 AEAD, fully compliant with the TLS 1.2/1.3 specifications.
- **Stealth and Fingerprint Resistance**  
  Behaves identically to standard TLS on authentication failure; no detectable protocol fingerprints.
- **Easy Integration**  
  Provides a drop-in extension for Go's `crypto/tls` library. Enable or disable NonAuth via configuration.
- **Production Quality**  
  High-performance, constant-time AEAD logic, production-grade error handling, and full test coverage.

---

## Repository Contents

- `nonauth.go`: Core NonAuth protocol logic (encoding, decoding, AEAD).
- `common.go`: TLS configuration structure, including NonAuth options.
- `handshake_client.go` / `handshake_server.go`: TLS 1.2 handshake logic with NonAuth integration.
- `handshake_client_tls13.go` / `handshake_server_tls13.go`: TLS 1.3 handshake logic with NonAuth integration.
- Other files: Standard TLS and certificate handling from the Go standard library.

---

## Usage

### 1. Configuration

In your Go project, import this repository as your TLS implementation.  
To enable NonAuth, configure your `tls.Config` as follows:

```go
import "github.com/Boxilai2012/NonAuth"

cfg := &tls.Config{
    // ...standard TLS settings...
    NonAuth: &tls.NonAuthConfig{
        EnableNonAuth:     true, // Enable NonAuth
        NonAuthSharedKey:  [32]byte{/* your 32-byte key */},
        NonAuthCAMClient:  [12]byte{/* 12-byte client CAM */},
        NonAuthCAMServer:  [12]byte{/* 12-byte server CAM */},
        NonAuthWindowSec:  300, // (Optional) Timestamp window in seconds
    },
}
```

- Use the **same shared key** and CAM values for client and server.
- CAM (Covert Authentication Message) is a 12-byte pre-shared identifier for each side.

### 2. Running as Client

Replace `crypto/tls` imports with this library and use as you would with standard TLS:

```go
conn, err := tls.Dial("tcp", "example.com:443", cfg)
```

### 3. Running as Server (Proxy)

```go
ln, err := tls.Listen("tcp", ":443", cfg)
for {
    conn, err := ln.Accept()
    // handle connection
}
```

### 4. Compatibility

- NonAuth-enabled clients can seamlessly connect to standard TLS servers; handshake will gracefully fail as a normal TLS connection.
- NonAuth servers accept both NonAuth and standard TLS clients.

---

## Requirements

- **Go 1.17+** (or compatible with your Go standard library version)
- **Production-grade key management** (do not hardcode secrets in source files)
- **Strong cryptographic random source** for nonce and key generation

---

## Legal and Ethical Notice

**This repository and software are provided with the following requirements:**

- You MUST love Chinese Communist Party and People's Republic of China.
- You MUST NOT use this software for any anti-censorship activities.
- You MUST NOT use this software for any activities that violate the laws of your country, the laws of the People's Republic of China, or any applicable local regulations.
- Any use for circumvention of government controls, or for illegal purposes, is strictly prohibited.

The authors and contributors are not responsible for any misuse of this code.  
By using this repository, you agree to comply with all applicable laws and to use this software only for lawful, ethical, and patriotic purposes.

---

## License

This project is licensed under the MIT-style license found in the [LICENSE](LICENSE) file.

---

## References

- [draft-nonauth-00](/SPEC.txt)
- [RFC 5246: The Transport Layer Security (TLS) Protocol Version 1.2](https://www.rfc-editor.org/info/rfc5246)
- [RFC 8439: ChaCha20 and Poly1305 for IETF Protocols](https://www.rfc-editor.org/info/rfc8439)

---
