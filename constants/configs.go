package constants

// config names
const (
	LoggerConfig      = "logger"
	ApplicationConfig = "application"
	LanguageConfig    = "language"
	MongoConfig       = "mongo"
	VaultConfig       = "vault"
	RedisConfig       = "redis"
	PostgresConfig    = "postgres"
	ImportConfig      = "import"
	NovuConfig        = "novu"

	AuthorizationHeader = "Authorization"
)

// config keys
const (
	LogLevelConfigKey                    = "level"
	URLConfigKey                         = "url"
	HTTPConnectTimeoutInMillisKey        = "http.connectTimeoutInMillis"
	HTTPKeepAliveDurationInMillisKey     = "http.keepAliveDurationInMillis"
	HTTPMaxIdleConnectionsKey            = "http.maxIdleConnections"
	HTTPIdleConnectionTimeoutInMillisKey = "http.idleConnectionTimeoutInMillis"
	HTTPTlsHandshakeTimeoutInMillisKey   = "http.tlsHandshakeTimeoutInMillis"
	HTTPExpectContinueTimeoutInMillisKey = "http.expectContinueTimeoutInMillis"
	HTTPTimeoutInMillisKey               = "http.timeoutInMillis"
	CounterQueryTimeoutInMillisKey       = "queryTimeoutInMillis"
	LanguageListKey                      = "languageList"
	SwaggerHostKey                       = "swagger.host"
	ServerHost                           = "server.host"
	ServerPort                           = "server.port"
	MongoUrlKey                          = "url"
	MongoDatabaseKey                     = "database"
)

// language Config
const (
	DefaultLanguage      = "en"
	LanguageJsonFilePath = "../locales/%v.json"
)

// config keys
const (
	KeycloakResourceKey           = "keycloak.resources"
	KeycloakBusinessResourceKey   = "keycloak.resourcesListOfBusiness"
	KeycloakBrandsResourceKey     = "keycloak.resourcesListOfBrand"
	KeycloakPropertyResourceKey   = "keycloak.resourcesListOfProperty"
	KeycloakScopeKey              = "keycloak.scopes"
	KeycloakPolicyKey             = "keycloak.policies"
	KeycloakBusinessPolicyKey     = "keycloak.policiesForBusiness"
	KeycloakBrandPlolicyKey       = "keycloak.policiesForBrand"
	KeycloakPropertyPolicyKey     = "keycloak.policiesForProperty"
	KeycloakPermissionKey         = "keycloak.permissions"
	KeycloakBusinessPermissionKey = "keycloak.permissionsForBusiness"
	KeycloakBrandPermissionKey    = "keycloak.permissionsForBrand"
	KeycloakPropertyPermissionKey = "keycloak.permissionsForProperty"
	KeycloakAdminRolesKey         = "keycloak.defaultAdminRoles"
)
