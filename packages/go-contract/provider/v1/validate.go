package providerv1

import (
	"fmt"
	"strings"

	apiprotocolv1 "code-code.internal/go-contract/api_protocol/v1"
	credentialv1 "code-code.internal/go-contract/credential/v1"
	modelv1 "code-code.internal/go-contract/model/v1"
)

// ValidateProviderSurface validates one provider surface metadata object.
func ValidateProviderSurface(surface *ProviderSurface) error {
	if surface == nil {
		return fmt.Errorf("providerv1: provider surface is nil")
	}
	if strings.TrimSpace(surface.GetSurfaceId()) == "" {
		return fmt.Errorf("providerv1: provider surface id is empty")
	}
	if len(surface.GetSupportedCredentialKinds()) == 0 {
		return fmt.Errorf("providerv1: supported credential kinds are required")
	}
	for _, kind := range surface.GetSupportedCredentialKinds() {
		if kind == credentialv1.CredentialKind_CREDENTIAL_KIND_UNSPECIFIED {
			return fmt.Errorf("providerv1: supported credential kind is unspecified")
		}
	}
	if surface.GetKind() == ProviderSurfaceKind_PROVIDER_SURFACE_KIND_UNSPECIFIED {
		return fmt.Errorf("providerv1: provider surface kind is unspecified")
	}
	if surface.GetKind() == ProviderSurfaceKind_PROVIDER_SURFACE_KIND_API {
		if len(surface.GetApi().GetSupportedProtocols()) == 0 {
			return fmt.Errorf("providerv1: api provider surface protocols are required")
		}
		for _, protocol := range surface.GetApi().GetSupportedProtocols() {
			if protocol == apiprotocolv1.Protocol_PROTOCOL_UNSPECIFIED {
				return fmt.Errorf("providerv1: supported protocol is unspecified")
			}
		}
	}
	return nil
}

// ValidateProvider validates one user-configured provider aggregate.
func ValidateProvider(provider *Provider) error {
	if provider == nil {
		return fmt.Errorf("providerv1: provider is nil")
	}
	if strings.TrimSpace(provider.GetProviderId()) == "" {
		return fmt.Errorf("providerv1: provider id is empty")
	}
	if strings.TrimSpace(provider.GetDisplayName()) == "" {
		return fmt.Errorf("providerv1: provider display name is empty")
	}
	if strings.TrimSpace(provider.GetSurfaceId()) == "" {
		return fmt.Errorf("providerv1: provider surface id is empty")
	}
	if provider.GetProviderCredentialRef() != nil {
		if err := ValidateProviderCredentialRef(provider.GetProviderCredentialRef()); err != nil {
			return err
		}
	}
	if err := ValidateProviderModels(provider.GetModels()); err != nil {
		return err
	}
	if custom := provider.GetCustomApiKeySurface(); custom != nil {
		if strings.TrimSpace(custom.GetBaseUrl()) == "" {
			return fmt.Errorf("providerv1: custom api key surface base_url is empty")
		}
		if custom.GetProtocol() == apiprotocolv1.Protocol_PROTOCOL_UNSPECIFIED {
			return fmt.Errorf("providerv1: custom api key surface protocol is unspecified")
		}
	}
	return nil
}

// ValidateCredentialCompatibility validates that one credential definition is compatible with one provider surface.
func ValidateCredentialCompatibility(surface *ProviderSurface, credential *credentialv1.CredentialDefinition) error {
	if err := ValidateProviderSurface(surface); err != nil {
		return err
	}
	if err := credentialv1.ValidateDefinition(credential); err != nil {
		return fmt.Errorf("providerv1: invalid credential definition: %w", err)
	}
	if !containsCredentialKind(surface.GetSupportedCredentialKinds(), credential.GetKind()) {
		return fmt.Errorf("providerv1: credential kind %s is not supported by surface %q", credential.GetKind().String(), surface.GetSurfaceId())
	}
	return nil
}

// ValidateProviderCredentialRef validates one provider credential reference.
func ValidateProviderCredentialRef(ref *ProviderCredentialRef) error {
	if ref == nil {
		return fmt.Errorf("providerv1: provider credential ref is nil")
	}
	if strings.TrimSpace(ref.GetProviderCredentialId()) == "" {
		return fmt.Errorf("providerv1: provider credential id is empty")
	}
	return nil
}

// ValidateProviderModelCatalog validates one provider-level model catalog.
func ValidateProviderModelCatalog(catalog *ProviderModelCatalog) error {
	if catalog == nil {
		return fmt.Errorf("providerv1: provider model catalog is nil")
	}
	if catalog.GetSource() == CatalogSource_CATALOG_SOURCE_UNSPECIFIED {
		return fmt.Errorf("providerv1: provider model catalog source is unspecified")
	}
	seen := map[string]struct{}{}
	for _, model := range catalog.GetModels() {
		if err := validateProviderModelCatalogEntry(model); err != nil {
			return err
		}
		providerModelID := strings.TrimSpace(model.GetProviderModelId())
		if _, ok := seen[providerModelID]; ok {
			return fmt.Errorf("providerv1: duplicate provider model id %q", providerModelID)
		}
		seen[providerModelID] = struct{}{}
	}
	return nil
}

// ValidateProviderModels validates provider-callable model entries.
func ValidateProviderModels(models []*ProviderModel) error {
	seen := map[string]struct{}{}
	for _, model := range models {
		if model == nil {
			return fmt.Errorf("providerv1: provider model is nil")
		}
		providerModelID := strings.TrimSpace(model.GetProviderModelId())
		if providerModelID == "" {
			return fmt.Errorf("providerv1: provider model id is empty")
		}
		if _, ok := seen[providerModelID]; ok {
			return fmt.Errorf("providerv1: duplicate provider model id %q", providerModelID)
		}
		seen[providerModelID] = struct{}{}
	}
	return nil
}

// ValidateResolvedProviderModel validates the final provider-routed model.
func ValidateResolvedProviderModel(surfaceMeta *ProviderSurface, provider *Provider, resolved *ResolvedProviderModel) error {
	if surfaceMeta != nil {
		if err := ValidateProviderSurface(surfaceMeta); err != nil {
			return err
		}
	}
	if provider != nil {
		if err := ValidateProvider(provider); err != nil {
			return err
		}
	}
	if resolved == nil {
		return fmt.Errorf("providerv1: resolved provider model is nil")
	}
	if strings.TrimSpace(resolved.GetSurfaceId()) == "" {
		return fmt.Errorf("providerv1: resolved provider model surface id is empty")
	}
	if provider != nil && resolved.GetSurfaceId() != provider.GetSurfaceId() {
		return fmt.Errorf("providerv1: resolved provider model surface id %q does not match provider surface id %q", resolved.GetSurfaceId(), provider.GetSurfaceId())
	}
	if strings.TrimSpace(resolved.GetProviderModelId()) == "" {
		return fmt.Errorf("providerv1: resolved provider model id is empty")
	}
	if resolved.GetProtocol() == apiprotocolv1.Protocol_PROTOCOL_UNSPECIFIED {
		return fmt.Errorf("providerv1: resolved provider model protocol is unspecified")
	}
	if provider != nil {
		if api := provider.GetRuntime().GetApi(); api != nil {
			if resolved.GetProtocol() != api.GetProtocol() {
				return fmt.Errorf("providerv1: resolved provider model protocol %s does not match surface protocol %s", resolved.GetProtocol().String(), api.GetProtocol().String())
			}
			if strings.TrimSpace(resolved.GetBaseUrl()) == "" {
				return fmt.Errorf("providerv1: resolved provider model requires base_url for api surface")
			}
		}
	}
	if err := ValidateResolvedProviderSurface(resolved.GetSurface()); err != nil {
		return fmt.Errorf("providerv1: invalid resolved provider surface: %w", err)
	}
	if resolved.GetSource() == CatalogSource_CATALOG_SOURCE_UNSPECIFIED {
		return fmt.Errorf("providerv1: resolved provider model source is unspecified")
	}
	if err := modelv1.ValidateResolvedModel(resolved.GetModel()); err != nil {
		return fmt.Errorf("providerv1: resolved provider model has invalid resolved model: %w", err)
	}
	return nil
}

// ValidateProviderSurfaceRuntime validates one callable provider surface runtime.
func ValidateProviderSurfaceRuntime(surface *ProviderSurfaceRuntime) error {
	if surface == nil {
		return fmt.Errorf("providerv1: provider surface runtime is nil")
	}
	if strings.TrimSpace(surface.GetDisplayName()) == "" {
		return fmt.Errorf("providerv1: provider surface display name is empty")
	}
	if surface.GetOrigin() == ProviderSurfaceOrigin_PROVIDER_SURFACE_ORIGIN_UNSPECIFIED {
		return fmt.Errorf("providerv1: provider surface origin is unspecified")
	}
	switch access := surface.GetAccess().(type) {
	case *ProviderSurfaceRuntime_Api:
		if access.Api.GetProtocol() == apiprotocolv1.Protocol_PROTOCOL_UNSPECIFIED {
			return fmt.Errorf("providerv1: api provider surface protocol is unspecified")
		}
		if strings.TrimSpace(access.Api.GetBaseUrl()) == "" {
			return fmt.Errorf("providerv1: api provider surface base_url is empty")
		}
	case *ProviderSurfaceRuntime_Cli:
		if strings.TrimSpace(access.Cli.GetCliId()) == "" {
			return fmt.Errorf("providerv1: cli provider surface cli_id is empty")
		}
	default:
		return fmt.Errorf("providerv1: provider surface access is unspecified")
	}
	return nil
}

// ValidateResolvedProviderSurface validates one final callable surface.
func ValidateResolvedProviderSurface(surface *ResolvedProviderSurface) error {
	if surface == nil {
		return fmt.Errorf("providerv1: resolved provider surface is nil")
	}
	if err := ValidateProviderSurfaceRuntime(surface.GetSurface()); err != nil {
		return err
	}
	if surface.GetAuth() != nil {
		if err := credentialv1.ValidateResolvedCredential(surface.GetAuth()); err != nil {
			return fmt.Errorf("providerv1: resolved provider surface auth is invalid: %w", err)
		}
	}
	return nil
}

func containsCredentialKind(values []credentialv1.CredentialKind, target credentialv1.CredentialKind) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func validateProviderModelCatalogEntry(model *ProviderModelCatalogEntry) error {
	if model == nil {
		return fmt.Errorf("providerv1: provider model entry is nil")
	}
	if strings.TrimSpace(model.GetProviderModelId()) == "" {
		return fmt.Errorf("providerv1: provider model id is empty")
	}
	return nil
}
