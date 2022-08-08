// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/authors": {
            "get": {
                "description": "Fetch Author",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Fetch Author",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "post": {
                "description": "Create Author",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Create Author",
                "parameters": [
                    {
                        "description": "Author to create",
                        "name": "author",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateAuthorReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/api/v1/authors/{id}": {
            "get": {
                "description": "Get Author",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Get Author",
                "parameters": [
                    {
                        "type": "string",
                        "description": "author id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "put": {
                "description": "Update Author",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Update Author",
                "parameters": [
                    {
                        "type": "string",
                        "description": "author id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Author to update",
                        "name": "author",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateAuthorReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "Delete Author",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Delete Author",
                "parameters": [
                    {
                        "type": "string",
                        "description": "author id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "request.CreateAuthorReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "request.UpdateAuthorReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.4",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Boilerplate",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}