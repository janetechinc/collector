package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"golang.org/x/net/http/httpproxy"

	"github.com/pganalyze/collector/util"
)

const DefaultAPIBaseURL = "https://api.pganalyze.com"

func getDefaultConfig() *ServerConfig {
	config := &ServerConfig{
		APIBaseURL:              DefaultAPIBaseURL,
		SectionName:             "default",
		QueryStatsInterval:      60,
		MaxCollectorConnections: 10,
	}

	// The environment variables are the default way to configure when running inside a Docker container.
	if apiKey := os.Getenv("PGA_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}
	if apiBaseURL := os.Getenv("PGA_API_BASEURL"); apiBaseURL != "" {
		config.APIBaseURL = apiBaseURL
	}
	if systemID := os.Getenv("PGA_API_SYSTEM_ID"); systemID != "" {
		config.SystemID = systemID
	}
	if systemType := os.Getenv("PGA_API_SYSTEM_TYPE"); systemType != "" {
		config.SystemType = systemType
	}
	if systemScope := os.Getenv("PGA_API_SYSTEM_SCOPE"); systemScope != "" {
		config.SystemScope = systemScope
	}
	if systemScopeFallback := os.Getenv("PGA_API_SYSTEM_SCOPE_FALLBACK"); systemScopeFallback != "" {
		config.SystemScopeFallback = systemScopeFallback
	}
	if enableReports := os.Getenv("PGA_ENABLE_REPORTS"); enableReports != "" && enableReports != "0" {
		config.EnableReports = true
	}
	if disableLogs := os.Getenv("PGA_DISABLE_LOGS"); disableLogs != "" && disableLogs != "0" {
		config.DisableLogs = true
	}
	if disableActivity := os.Getenv("PGA_DISABLE_ACTIVITY"); disableActivity != "" && disableActivity != "0" {
		config.DisableActivity = true
	}
	if enableLogExplain := os.Getenv("PGA_ENABLE_LOG_EXPLAIN"); enableLogExplain != "" && enableLogExplain != "0" {
		config.EnableLogExplain = true
	}
	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		config.DbURL = dbURL
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.DbName = dbName
	}
	if dbAllNames := os.Getenv("DB_ALL_NAMES"); dbAllNames == "1" {
		config.DbAllNames = true
	}
	if dbUsername := os.Getenv("DB_USERNAME"); dbUsername != "" {
		config.DbUsername = dbUsername
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.DbPassword = dbPassword
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.DbHost = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.DbPort, _ = strconv.Atoi(dbPort)
	}
	if dbSslMode := os.Getenv("DB_SSLMODE"); dbSslMode != "" {
		config.DbSslMode = dbSslMode
	}
	if dbSslRootCert := os.Getenv("DB_SSLROOTCERT"); dbSslRootCert != "" {
		config.DbSslRootCert = dbSslRootCert
	}
	if dbSslRootCertContents := os.Getenv("DB_SSLROOTCERT_CONTENTS"); dbSslRootCertContents != "" {
		config.DbSslRootCertContents = dbSslRootCertContents
	}
	if dbSslCert := os.Getenv("DB_SSLCERT"); dbSslCert != "" {
		config.DbSslCert = dbSslCert
	}
	if dbSslCertContents := os.Getenv("DB_SSLCERT_CONTENTS"); dbSslCertContents != "" {
		config.DbSslCertContents = dbSslCertContents
	}
	if dbSslKey := os.Getenv("DB_SSLKEY"); dbSslKey != "" {
		config.DbSslKey = dbSslKey
	}
	if dbSslKeyContents := os.Getenv("DB_SSLKEY_CONTENTS"); dbSslKeyContents != "" {
		config.DbSslKeyContents = dbSslKeyContents
	}
	if awsRegion := os.Getenv("AWS_REGION"); awsRegion != "" {
		config.AwsRegion = awsRegion
	}
	if awsAccountID := os.Getenv("AWS_ACCOUNT_ID"); awsAccountID != "" {
		config.AwsAccountID = awsAccountID
	}
	if awsInstanceID := os.Getenv("AWS_INSTANCE_ID"); awsInstanceID != "" {
		config.AwsDbInstanceID = awsInstanceID
	}
	if awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID"); awsAccessKeyID != "" {
		config.AwsAccessKeyID = awsAccessKeyID
	}
	if awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY"); awsSecretAccessKey != "" {
		config.AwsSecretAccessKey = awsSecretAccessKey
	}
	if awsAssumeRole := os.Getenv("AWS_ASSUME_ROLE"); awsAssumeRole != "" {
		config.AwsAssumeRole = awsAssumeRole
	}
	if awsWebIdentityTokenFile := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE"); awsWebIdentityTokenFile != "" {
		config.AwsWebIdentityTokenFile = awsWebIdentityTokenFile
	}
	if awsRoleArn := os.Getenv("AWS_ROLE_ARN"); awsRoleArn != "" {
		config.AwsRoleArn = awsRoleArn
	}
	if awsEndpointSigningRegion := os.Getenv("AWS_ENDPOINT_SIGNING_REGION"); awsEndpointSigningRegion != "" {
		config.AwsEndpointSigningRegion = awsEndpointSigningRegion
	}
	if awsEndpointRdsURL := os.Getenv("AWS_ENDPOINT_RDS_URL"); awsEndpointRdsURL != "" {
		config.AwsEndpointRdsURL = awsEndpointRdsURL
	}
	if awsEndpointEc2URL := os.Getenv("AWS_ENDPOINT_EC2_URL"); awsEndpointEc2URL != "" {
		config.AwsEndpointEc2URL = awsEndpointEc2URL
	}
	if awsEndpointCloudwatchURL := os.Getenv("AWS_ENDPOINT_CLOUDWATCH_URL"); awsEndpointCloudwatchURL != "" {
		config.AwsEndpointCloudwatchURL = awsEndpointCloudwatchURL
	}
	if awsEndpointCloudwatchLogsURL := os.Getenv("AWS_ENDPOINT_CLOUDWATCH_LOGS_URL"); awsEndpointCloudwatchLogsURL != "" {
		config.AwsEndpointCloudwatchLogsURL = awsEndpointCloudwatchLogsURL
	}
	if azureDbServerName := os.Getenv("AZURE_DB_SERVER_NAME"); azureDbServerName != "" {
		config.AzureDbServerName = azureDbServerName
	}
	if azureEventhubNamespace := os.Getenv("AZURE_EVENTHUB_NAMESPACE"); azureEventhubNamespace != "" {
		config.AzureEventhubNamespace = azureEventhubNamespace
	}
	if azureEventhubName := os.Getenv("AZURE_EVENTHUB_NAME"); azureEventhubName != "" {
		config.AzureEventhubName = azureEventhubName
	}
	if azureADTenantID := os.Getenv("AZURE_AD_TENANT_ID"); azureADTenantID != "" {
		config.AzureADTenantID = azureADTenantID
	}
	if azureADClientID := os.Getenv("AZURE_AD_CLIENT_ID"); azureADClientID != "" {
		config.AzureADClientID = azureADClientID
	}
	if azureADClientSecret := os.Getenv("AZURE_AD_CLIENT_SECRET"); azureADClientSecret != "" {
		config.AzureADClientSecret = azureADClientSecret
	}
	if azureADCertificatePath := os.Getenv("AZURE_AD_CERTIFICATE_PATH"); azureADCertificatePath != "" {
		config.AzureADCertificatePath = azureADCertificatePath
	}
	if azureADCertificatePassword := os.Getenv("AZURE_AD_CERTIFICATE_PASSWORD"); azureADCertificatePassword != "" {
		config.AzureADCertificatePassword = azureADCertificatePassword
	}
	if gcpCloudSQLInstanceID := os.Getenv("GCP_CLOUDSQL_INSTANCE_ID"); gcpCloudSQLInstanceID != "" {
		config.GcpCloudSQLInstanceID = gcpCloudSQLInstanceID
	}
	if gcpPubsubSubscription := os.Getenv("GCP_PUBSUB_SUBSCRIPTION"); gcpPubsubSubscription != "" {
		config.GcpPubsubSubscription = gcpPubsubSubscription
	}
	if gcpCredentialsFile := os.Getenv("GCP_CREDENTIALS_FILE"); gcpCredentialsFile != "" {
		config.GcpCredentialsFile = gcpCredentialsFile
	}
	if gcpProjectID := os.Getenv("GCP_PROJECT_ID"); gcpProjectID != "" {
		config.GcpProjectID = gcpProjectID
	}
	if logLocation := os.Getenv("LOG_LOCATION"); logLocation != "" {
		config.LogLocation = logLocation
	}
	if logSyslogServer := os.Getenv("LOG_SYSLOG_SERVER"); logSyslogServer != "" {
		config.LogSyslogServer = logSyslogServer
	}
	// Note: We don't support LogDockerTail here since it would require the "docker"
	// binary inside the pganalyze container (as well as full Docker access), instead
	// the approach for using pganalyze as a sidecar container alongside Postgres
	// currently requires writing to a file and then mounting that as a volume
	// inside the pganalyze container.
	if ignoreTablePattern := os.Getenv("IGNORE_TABLE_PATTERN"); ignoreTablePattern != "" {
		config.IgnoreTablePattern = ignoreTablePattern
	}
	if ignoreSchemaRegexp := os.Getenv("IGNORE_SCHEMA_REGEXP"); ignoreSchemaRegexp != "" {
		config.IgnoreSchemaRegexp = ignoreSchemaRegexp
	}
	if queryStatsInterval := os.Getenv("QUERY_STATS_INTERVAL"); queryStatsInterval != "" {
		config.QueryStatsInterval, _ = strconv.Atoi(queryStatsInterval)
	}
	if maxCollectorConnections := os.Getenv("MAX_COLLECTOR_CONNECTION"); maxCollectorConnections != "" {
		config.MaxCollectorConnections, _ = strconv.Atoi(maxCollectorConnections)
	}
	if skipIfReplica := os.Getenv("SKIP_IF_REPLICA"); skipIfReplica != "" && skipIfReplica != "0" {
		config.SkipIfReplica = true
	}
	if filterLogSecret := os.Getenv("FILTER_LOG_SECRET"); filterLogSecret != "" {
		config.FilterLogSecret = filterLogSecret
	}
	if filterQuerySample := os.Getenv("FILTER_QUERY_SAMPLE"); filterQuerySample != "" {
		config.FilterQuerySample = filterQuerySample
	}
	if filterQueryText := os.Getenv("FILTER_QUERY_TEXT"); filterQueryText != "" {
		config.FilterQueryText = filterQueryText
	}
	if httpProxy := os.Getenv("HTTP_PROXY"); httpProxy != "" {
		config.HTTPProxy = httpProxy
	}
	if httpProxy := os.Getenv("http_proxy"); httpProxy != "" {
		config.HTTPProxy = httpProxy
	}
	if httpsProxy := os.Getenv("HTTPS_PROXY"); httpsProxy != "" {
		config.HTTPSProxy = httpsProxy
	}
	if httpsProxy := os.Getenv("https_proxy"); httpsProxy != "" {
		config.HTTPSProxy = httpsProxy
	}
	if noProxy := os.Getenv("NO_PROXY"); noProxy != "" {
		config.NoProxy = noProxy
	}
	if noProxy := os.Getenv("no_proxy"); noProxy != "" {
		config.NoProxy = noProxy
	}

	return config
}

func CreateHTTPClient(conf ServerConfig) *http.Client {
	requireSSL := conf.APIBaseURL == DefaultAPIBaseURL
	proxyConfig := httpproxy.Config{
		HTTPProxy:  conf.HTTPProxy,
		HTTPSProxy: conf.HTTPSProxy,
		NoProxy:    conf.NoProxy,
	}
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return proxyConfig.ProxyFunc()(req.URL)
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if requireSSL {
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			matchesProxyURL := false
			if proxyConfig.HTTPProxy != "" {
				proxyURL, err := url.Parse(proxyConfig.HTTPProxy)
				if err == nil && proxyURL.Host == addr {
					matchesProxyURL = true
				}
			}
			if proxyConfig.HTTPSProxy != "" {
				proxyURL, err := url.Parse(proxyConfig.HTTPSProxy)
				if err == nil && proxyURL.Host == addr {
					matchesProxyURL = true
				}
			}
			// Require secure conection for everything except proxies
			if !matchesProxyURL && !strings.HasSuffix(addr, ":443") {
				return nil, fmt.Errorf("Unencrypted connection is not permitted by pganalyze configuration")
			}
			return (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second, DualStack: true}).DialContext(ctx, network, addr)
		}
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	return &http.Client{
		Timeout:   120 * time.Second,
		Transport: transport,
	}
}

// CreateEC2IMDSHTTPClient - Create HTTP client for EC2 instance meta data service (IMDS)
func CreateEC2IMDSHTTPClient(conf ServerConfig) *http.Client {
	// Match https://github.com/aws/aws-sdk-go/pull/3066
	return &http.Client{
		Timeout: 1 * time.Second,
	}
}

func writeValueToTempfile(value string) (string, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(value)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func preprocessConfig(config *ServerConfig) (*ServerConfig, error) {
	var err error

	host := config.GetDbHost()
	if strings.HasSuffix(host, ".rds.amazonaws.com") {
		parts := strings.SplitN(host, ".", 4)
		if len(parts) == 4 && parts[3] == "rds.amazonaws.com" { // Safety check for any escaping issues
			if config.AwsDbInstanceID == "" {
				config.AwsDbInstanceID = parts[0]
			}
			if config.AwsAccountID == "" {
				config.AwsAccountID = parts[1]
			}
			if config.AwsRegion == "" {
				config.AwsRegion = parts[2]
			}
		}
	} else if strings.HasSuffix(host, ".postgres.database.azure.com") {
		parts := strings.SplitN(host, ".", 2)
		if len(parts) == 2 && parts[1] == "postgres.database.azure.com" { // Safety check for any escaping issues
			if config.AzureDbServerName == "" {
				config.AzureDbServerName = parts[0]
			}
		}
	}

	// This is primarily for backwards compatibility when using the IP address of an instance
	// combined with only specifying its name, but not its region.
	if config.AwsDbInstanceID != "" && config.AwsRegion == "" {
		config.AwsRegion = "us-east-1"
	}

	if config.GcpCloudSQLInstanceID != "" && strings.Count(config.GcpCloudSQLInstanceID, ":") == 2 {
		instanceParts := strings.SplitN(config.GcpCloudSQLInstanceID, ":", 3)
		config.GcpProjectID = instanceParts[0]
		config.GcpCloudSQLInstanceID = instanceParts[2]
	}

	dbNameParts := []string{}
	for _, s := range strings.Split(config.DbName, ",") {
		dbNameParts = append(dbNameParts, strings.TrimSpace(s))
	}
	config.DbName = dbNameParts[0]
	if len(dbNameParts) == 2 && dbNameParts[1] == "*" {
		config.DbAllNames = true
	} else {
		config.DbExtraNames = dbNameParts[1:]
	}

	if config.DbSslRootCertContents != "" {
		config.DbSslRootCert, err = writeValueToTempfile(config.DbSslRootCertContents)
		if err != nil {
			return config, err
		}
	}

	if config.DbSslCertContents != "" {
		config.DbSslCert, err = writeValueToTempfile(config.DbSslCertContents)
	}

	if config.DbSslKeyContents != "" {
		config.DbSslKey, err = writeValueToTempfile(config.DbSslKeyContents)
	}

	if config.AwsEndpointSigningRegionLegacy != "" && config.AwsEndpointSigningRegion == "" {
		config.AwsEndpointSigningRegion = config.AwsEndpointSigningRegionLegacy
	}

	return config, nil
}

// Read - Reads the configuration from the specified filename, or fall back to the default config
func Read(logger *util.Logger, filename string) (Config, error) {
	var conf Config
	var err error

	if _, err = os.Stat(filename); err == nil {
		configFile, err := ini.LoadSources(ini.LoadOptions{SpaceBeforeInlineComment: true}, filename)
		if err != nil {
			return conf, err
		}

		defaultConfig := getDefaultConfig()

		pgaSection, err := configFile.GetSection("pganalyze")
		if err != nil {
			return conf, fmt.Errorf("Failed to find [pganalyze] section in config: %s", err)
		}
		err = pgaSection.MapTo(defaultConfig)
		if err != nil {
			return conf, fmt.Errorf("Failed to map [pganalyze] section in config: %s", err)
		}

		sections := configFile.Sections()
		for _, section := range sections {
			config := &ServerConfig{}
			*config = *defaultConfig

			err = section.MapTo(config)
			if err != nil {
				return conf, err
			}

			config, err = preprocessConfig(config)
			if err != nil {
				return conf, err
			}
			config.SectionName = section.Name()
			config.SystemType, config.SystemScope, config.SystemScopeFallback, config.SystemID = identifySystem(*config)

			config.Identifier = ServerIdentifier{
				APIKey:      config.APIKey,
				APIBaseURL:  config.APIBaseURL,
				SystemID:    config.SystemID,
				SystemType:  config.SystemType,
				SystemScope: config.SystemScope,
			}

			if config.GetDbName() != "" {
				// Ensure we have no duplicate identifiers within one collector
				skip := false
				for _, server := range conf.Servers {
					if config.Identifier == server.Identifier {
						skip = true
					}
				}
				if skip {
					logger.PrintError("Skipping config section %s, detected as duplicate", config.SectionName)
				} else {
					conf.Servers = append(conf.Servers, *config)
				}
			}

			if config.DbURL != "" {
				_, err := url.Parse(config.DbURL)
				if err != nil {
					prefixedLogger := logger.WithPrefix(config.SectionName)
					prefixedLogger.PrintError("Could not parse db_url; check URL format and note that any special characters must be percent-encoded")
				}
			}
		}

		if len(conf.Servers) == 0 {
			return conf, fmt.Errorf("Configuration file is empty, please edit %s and reload the collector", filename)
		}
	} else {
		if os.Getenv("DYNO") != "" && os.Getenv("PORT") != "" {
			for _, kv := range os.Environ() {
				parts := strings.Split(kv, "=")
				if strings.HasSuffix(parts[0], "_URL") {
					config := getDefaultConfig()
					config, err = preprocessConfig(config)
					if err != nil {
						return conf, err
					}
					config.SectionName = parts[0]
					config.SystemID = strings.Replace(parts[0], "_URL", "", 1)
					config.SystemType = "heroku"
					config.DbURL = parts[1]
					conf.Servers = append(conf.Servers, *config)
				}
			}
		} else if os.Getenv("PGA_API_KEY") != "" {
			config := getDefaultConfig()
			config, err = preprocessConfig(config)
			if err != nil {
				return conf, err
			}
			config.SystemType, config.SystemScope, config.SystemScopeFallback, config.SystemID = identifySystem(*config)
			conf.Servers = append(conf.Servers, *config)
		} else {
			return conf, fmt.Errorf("No configuration file found at %s, and no environment variables set", filename)
		}
	}

	var hasIgnoreTablePattern = false
	for _, server := range conf.Servers {
		if server.IgnoreTablePattern != "" {
			hasIgnoreTablePattern = true
			break
		}
	}

	if hasIgnoreTablePattern {
		if os.Getenv("IGNORE_TABLE_PATTERN") != "" {
			logger.PrintVerbose("Deprecated: Setting IGNORE_TABLE_PATTERN is deprecated; please use IGNORE_SCHEMA_REGEXP instead")
		} else {
			logger.PrintVerbose("Deprecated: Setting ignore_table_pattern is deprecated; please use ignore_schema_regexp instead")
		}
	}

	return conf, nil
}
