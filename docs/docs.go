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
        "/api/v1/auth/signin": {
            "post": {
                "description": "SignIn",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignIn",
                "parameters": [
                    {
                        "description": "SignIn user",
                        "name": "signin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignInReq"
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
        "/api/v1/auth/signup": {
            "post": {
                "description": "SignUp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignUp",
                "parameters": [
                    {
                        "description": "SignUp user",
                        "name": "signup",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignUpReq"
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
        "/api/v1/todos": {
            "get": {
                "security": [
                    {
                        "JwtToken": []
                    }
                ],
                "description": "Fetch Todo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Fetch Todo",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "JwtToken": []
                    }
                ],
                "description": "Create Todo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Create Todo",
                "parameters": [
                    {
                        "description": "Todo to create",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateTodoReq"
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
        "/api/v1/todos/{id}": {
            "get": {
                "security": [
                    {
                        "JwtToken": []
                    }
                ],
                "description": "Get Todo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Get Todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "todo id",
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
                "security": [
                    {
                        "JwtToken": []
                    }
                ],
                "description": "Update Todo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Update Todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "todo id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Todo to update",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateTodoReq"
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
                "security": [
                    {
                        "JwtToken": []
                    }
                ],
                "description": "Delete Todo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Delete Todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "todo id",
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
        "request.CreateTodoReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "request.SignInReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.SignUpReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.UpdateTodoReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "JwtToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
