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
        "/backup/{project_id}/{report_id}": {
            "get": {
                "description": "Backup Report By ID",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Backup Report By ID",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "default",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "latest",
                        "description": "default",
                        "name": "report_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": true,
                        "description": "default",
                        "name": "should_delete",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "default": true,
                        "description": "default",
                        "name": "async",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/clean-history": {
            "get": {
                "description": "Clean history project",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Clean history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "projectId",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/clean-results": {
            "get": {
                "description": "Clean allure result files on server",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Clean results",
                "parameters": [
                    {
                        "type": "string",
                        "description": "projectId",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/config": {
            "get": {
                "description": "Get app config",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "General"
                ],
                "summary": "Get App Config",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/generate-report": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Generate report from sent results",
                "parameters": [
                    {
                        "type": "string",
                        "description": "projectId",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "executionName",
                        "name": "execution_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "executionName",
                        "name": "execution_from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "executionName",
                        "name": "execution_type",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "executionName",
                        "name": "backup_latest",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/projects": {
            "get": {
                "description": "Get All projects",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Get All projects",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            },
            "post": {
                "description": "Create Project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Create Project",
                "parameters": [
                    {
                        "description": "default",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateProjectRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/projects/batch-delete": {
            "post": {
                "description": "Batch Delete Project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Batch Delete Project",
                "parameters": [
                    {
                        "description": "default",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.BatchDeleteRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/projects/{project_id}": {
            "get": {
                "description": "Get Project By ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Get Project By ID",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "default",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            },
            "delete": {
                "description": "Delete Project",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Delete Project",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "default",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/projects/{project_id}/reports/{path}": {
            "get": {
                "description": "Get Path in a project directory",
                "consumes": [
                    "*/*"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Get Path in a project directory",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "default",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "latest/widgets/summary.json",
                        "description": "default",
                        "name": "path",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/report/download": {
            "get": {
                "description": "Download latest report csv",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Download latest report csv",
                "parameters": [
                    {
                        "type": "string",
                        "default": "default",
                        "description": "default",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/send-results": {
            "post": {
                "description": "Send allure result files to server",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Send results",
                "parameters": [
                    {
                        "type": "string",
                        "description": "projectId",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "create project if not exists",
                        "name": "force_project_creation",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "file"
                        },
                        "collectionFormat": "csv",
                        "description": "result files",
                        "name": "files[]",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.BatchDeleteRequest": {
            "type": "object",
            "properties": {
                "async": {
                    "type": "boolean"
                },
                "projects": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "test123",
                        "test456",
                        "test789"
                    ]
                }
            }
        },
        "models.CreateProjectRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "default"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
