# tra 🍃

[![release](https://img.shields.io/github/v/release/dundunlabs/tra)](https://github.com/dundunlabs/tra/releases)
[![ci](https://github.com/dundunlabs/tra/actions/workflows/ci.yml/badge.svg)](https://github.com/dundunlabs/tra/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/dundunlabs/tra/branch/main/graph/badge.svg?token=qQQkGXZRKC)](https://codecov.io/gh/dundunlabs/tra)
[![license](https://img.shields.io/github/license/dundunlabs/tra)](https://github.com/dundunlabs/tra/blob/main/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://makeapullrequest.com)

Simple and lightweight HTTP router for [Go](https://go.dev/) with some basic features:

- **Routing**: match any route like "static" (e.g. `/api/users`), "dynamic" (e.g. `/users/:id`) and "wildcard" (e.g. `/public/*`)
- **Middleware**: execute any code before and after route handler (e.g. authentication, logging, ...)
- **Error handling**: create, custom, return error without tears

## Installation 
```bash
go get -u github.com/dundunlabs/tra
```

## Example
https://github.com/dundunlabs/tra/tree/main/example

## License
[MIT](https://github.com/dundunlabs/tra/blob/main/LICENSE)

---
Enjoy your coding 😊 
