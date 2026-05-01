package providerv1

import (
	"fmt"
	"strings"

	apiprotocolv1 "code-code.internal/go-contract/api_protocol/v1"
	"google.golang.org/protobuf/proto"
)

func ValidateProviderEndpoint(endpoint *ProviderEndpoint) error {
	if endpoint == nil {
		return fmt.Errorf("providerv1: provider endpoint is nil")
	}
	switch endpoint.GetType() {
	case ProviderEndpointType_PROVIDER_ENDPOINT_TYPE_API:
		api := endpoint.GetApi()
		if api == nil {
			return fmt.Errorf("providerv1: api provider endpoint is empty")
		}
		if api.GetProtocol() == apiprotocolv1.Protocol_PROTOCOL_UNSPECIFIED {
			return fmt.Errorf("providerv1: api provider endpoint protocol is unspecified")
		}
		if strings.TrimSpace(api.GetBaseUrl()) == "" {
			return fmt.Errorf("providerv1: api provider endpoint base_url is empty")
		}
	case ProviderEndpointType_PROVIDER_ENDPOINT_TYPE_CLI:
		if strings.TrimSpace(endpoint.GetCli().GetCliId()) == "" {
			return fmt.Errorf("providerv1: cli provider endpoint cli_id is empty")
		}
	default:
		return fmt.Errorf("providerv1: provider endpoint type is unspecified")
	}
	return nil
}

func EndpointProtocol(endpoint *ProviderEndpoint) apiprotocolv1.Protocol {
	if endpoint == nil || endpoint.GetApi() == nil {
		return apiprotocolv1.Protocol_PROTOCOL_UNSPECIFIED
	}
	return endpoint.GetApi().GetProtocol()
}

func EndpointBaseURL(endpoint *ProviderEndpoint) string {
	if endpoint == nil || endpoint.GetApi() == nil {
		return ""
	}
	return strings.TrimSpace(endpoint.GetApi().GetBaseUrl())
}

func EndpointCLIID(endpoint *ProviderEndpoint) string {
	if endpoint == nil || endpoint.GetCli() == nil {
		return ""
	}
	return strings.TrimSpace(endpoint.GetCli().GetCliId())
}

func EndpointKey(endpoint *ProviderEndpoint) string {
	if endpoint == nil {
		return ""
	}
	switch endpoint.GetType() {
	case ProviderEndpointType_PROVIDER_ENDPOINT_TYPE_API:
		if api := endpoint.GetApi(); api != nil {
			return "api:" + api.GetProtocol().String() + ":" + strings.TrimSpace(api.GetBaseUrl())
		}
	case ProviderEndpointType_PROVIDER_ENDPOINT_TYPE_CLI:
		if cliID := EndpointCLIID(endpoint); cliID != "" {
			return "cli:" + cliID
		}
	}
	return ""
}

func CloneProviderEndpoint(endpoint *ProviderEndpoint) *ProviderEndpoint {
	if endpoint == nil {
		return nil
	}
	return proto.Clone(endpoint).(*ProviderEndpoint)
}
