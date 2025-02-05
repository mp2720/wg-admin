// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all registered users. Admin only.",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/v1.UserResponse"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Register a new user. Admin only.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Register a new user.",
                "parameters": [
                    {
                        "description": "user data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Ok",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/v1.UserResponse"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "409": {
                        "description": "User with given name or key already exists",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    }
                }
            }
        },
        "/users/{uuid}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get user by UUID. Admin only.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get user by UUID.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user's UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/v1.UserResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    }
                }
            }
        },
        "/users/{uuid}/token": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Issue token for the user invalidating the previous. All users can issue tokens for their accounts. Only admin can issue token for other users.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Issue token for the user.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user's UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/v1.UserResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/v1.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.APIError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.RegisterUserRequest": {
            "type": "object",
            "properties": {
                "fare": {
                    "type": "string"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "max_addresses": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "private_key": {
                    "type": "string"
                }
            }
        },
        "v1.UserResponse": {
            "type": "object",
            "properties": {
                "address_count": {
                    "type": "integer"
                },
                "fare": {
                    "type": "string"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "is_banned": {
                    "type": "boolean"
                },
                "last_seen_at": {
                    "type": "string"
                },
                "links": {
                    "$ref": "#/definitions/v1.UserResponseLinks"
                },
                "max_addresses": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "paid_by_time": {
                    "type": "string"
                },
                "public_key": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "v1.UserResponseLinks": {
            "type": "object",
            "properties": {
                "addresses": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "X-Token",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Wireguard admin server",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
