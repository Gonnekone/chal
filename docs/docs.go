// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "url": "https://t.me/Gonnekone",
            "email": "opacha2018@yandex.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Retrieves a list of cars based on specified filters.",
                "tags": [
                    "cars"
                ],
                "summary": "Get cars",
                "parameters": [
                    {
                        "description": "Request body containing filters",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_get.Request"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "client error"
                    }
                }
            },
            "post": {
                "description": "Saves a list of cars with the provided registration numbers.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Save cars",
                "parameters": [
                    {
                        "description": "Request body containing registration numbers",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_save.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "client error"
                    }
                }
            },
            "delete": {
                "description": "Deletes a car by its identifier.",
                "tags": [
                    "cars"
                ],
                "summary": "Delete a car",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Car identifier",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "client error"
                    }
                }
            },
            "patch": {
                "description": "Updates the details of a car with the provided ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Update car",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Car ID to update",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request body containing updated car details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_update.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "client error"
                    }
                }
            }
        }
    },
    "definitions": {
        "internal_http-server_handlers_get.Request": {
            "type": "object",
            "properties": {
                "mark": {
                    "type": "string",
                    "example": "BMW"
                },
                "model": {
                    "type": "string",
                    "example": "X5"
                },
                "ownerName": {
                    "type": "string",
                    "example": "Ivan"
                },
                "ownerPatronymic": {
                    "type": "string",
                    "example": "Ivanovich"
                },
                "ownerSurname": {
                    "type": "string",
                    "example": "Ivanov"
                },
                "regNum": {
                    "type": "string",
                    "example": "AAA111"
                },
                "year": {
                    "type": "integer",
                    "example": 2020
                }
            }
        },
        "internal_http-server_handlers_save.Request": {
            "type": "object",
            "properties": {
                "regNums": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "AAA111",
                        " BBB222",
                        " CCC333"
                    ]
                }
            }
        },
        "internal_http-server_handlers_update.Request": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "mark": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "regNum": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8082",
	BasePath:         "",
	Schemes:          []string{"http"},
	Title:            "Cars catalog API",
	Description:      "This is a testovoe zadanie.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
