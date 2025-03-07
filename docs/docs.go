// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag/v2"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "components": {"schemas":{"apperror.AppError":{"properties":{"code":{"type":"integer"},"err":{},"message":{"type":"string"}},"type":"object"},"domain.Book":{"properties":{"author":{"type":"string"},"id":{"type":"integer"},"title":{"type":"string"}},"type":"object"},"dto.Pagination":{"properties":{"current_page":{"type":"integer"},"last_page":{"type":"integer"},"limit":{"type":"integer"},"total":{"type":"integer"}},"type":"object"},"dto.PaginationResponse-domain_Book":{"properties":{"data":{"items":{"$ref":"#/components/schemas/domain.Book"},"type":"array","uniqueItems":false},"pagination":{"$ref":"#/components/schemas/dto.Pagination"}},"type":"object"},"dto.SuccessResponse-domain_Book":{"properties":{"data":{"$ref":"#/components/schemas/domain.Book"}},"type":"object"}},"securitySchemes":{"Bearer":{"description":"Bearer token authentication","in":"header","name":"Authorization","type":"apiKey"}}},
    "info": {"description":"{{escape .Description}}","title":"{{.Title}}","version":"{{.Version}}"},
    "externalDocs": {"description":"","url":""},
    "paths": {"/books":{"get":{"description":"get books","parameters":[{"description":"Number of history to be retrieved","in":"query","name":"limit","schema":{"type":"integer"}},{"description":"Page to retrieved","in":"query","name":"page","schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.PaginationResponse-domain_Book"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Bad Request"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Internal Server Error"}},"summary":"GetBooks","tags":["book"]},"post":{"description":"create book by title and author","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/domain.Book"}}},"description":"Book Data","required":true},"responses":{"201":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-domain_Book"}}},"description":"Created"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Bad Request"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Internal Server Error"}},"summary":"CreateBook","tags":["book"]}},"/books/{id}":{"delete":{"description":"update book data","parameters":[{"description":"Book ID","in":"path","name":"id","required":true,"schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-domain_Book"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Bad Request"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Internal Server Error"}},"summary":"UpdateBook","tags":["book"]},"get":{"description":"get book by id","parameters":[{"description":"Book ID","in":"path","name":"id","required":true,"schema":{"type":"integer"}}],"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-domain_Book"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Bad Request"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Internal Server Error"}},"summary":"GetBook","tags":["book"]},"patch":{"description":"update book data","parameters":[{"description":"Book ID","in":"path","name":"id","required":true,"schema":{"type":"integer"}}],"requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/domain.Book"}}},"description":"Book Data","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/dto.SuccessResponse-domain_Book"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Bad Request"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/apperror.AppError"}}},"description":"Internal Server Error"}},"summary":"UpdateBook","tags":["book"]}}},
    "openapi": "3.1.0"
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Title:            "GO-FIBER-TEMPLATE API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
