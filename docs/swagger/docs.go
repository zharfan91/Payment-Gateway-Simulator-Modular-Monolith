package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "description": "Modular monolith payment gateway simulator.",
    "title": "Payment Gateway Simulator API",
    "contact": {},
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/api/v1",
  "paths": {
    "/auth/register": {"post": {"summary": "Register user", "responses": {"201": {"description": "created"}}}},
    "/auth/login": {"post": {"summary": "Login user", "responses": {"200": {"description": "ok"}}}},
    "/payments": {"post": {"summary": "Create payment", "responses": {"201": {"description": "created"}}}},
    "/settlements": {"post": {"summary": "Create settlement", "responses": {"201": {"description": "created"}}}}
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Payment Gateway Simulator API",
	Description:      "Modular monolith payment gateway simulator.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
