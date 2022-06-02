{
    "openapi": "3.0.1",
    "info": {
    "title": "entando-template",
        "version": "1"
},
    "servers": [
    {
        "url": "{{.ServersURL}}",
        "description": "Generated server url"
    }
],
    "security": [
    {
        "keycloak": [
            "read",
            "write"
        ]
    }
],
    "components": {
    "schemas": {
        "Response": {
            "type": "object",
                "properties": {
                "metric": {
                    "type": "string"
                }
            }
        }
    },
    "securitySchemes": {
        "keycloak": {
            "type": "oauth2",
                "flows": {
                "implicit": {
                    "authorizationUrl": "{{.AuthorizationURL}}/realms/{{.KeycloakRealm}}/protocol/openid-connect/auth",
                        "scopes": {}
                }
            }
        }
    }
},
    "paths": {
    "/api/example": {
        "get": {
            "description": "Returns a custom metric",
                "responses": {
                "200": {
                    "description": "the metric"
                }
            }
        }
    }
},
    "tags": []
}
