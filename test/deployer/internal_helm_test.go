package deployer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	pkgdeployer "github.com/kgateway-dev/kgateway/v2/pkg/deployer"
	"github.com/kgateway-dev/kgateway/v2/pkg/schemes"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
)

func TestRenderHelmChart(t *testing.T) {
	tests := []HelmTestCase{
		{
			Name:      "basic gateway with default gatewayclass and no gwparams",
			InputFile: "base-gateway",
		},
		{
			Name:      "gateway with replicas GWP via GWC",
			InputFile: "gwc-with-replicas",
		},
		{
			Name:      "gwparams with omitDefaultSecurityContext via GWC",
			InputFile: "omit-default-security-context",
		},
		{
			Name:      "gwparams with omitDefaultSecurityContext via GW",
			InputFile: "omit-default-security-context-via-gw",
		},
		{
			Name:      "agentgateway",
			InputFile: "agentgateway",
		},
		{
			Name:      "agentgateway OmitDefaultSecurityContext true GWP via GWC",
			InputFile: "agentgateway-omitdefaultsecuritycontext",
		},
		{
			Name:      "agentgateway OmitDefaultSecurityContext true GWP via GW",
			InputFile: "agentgateway-omitdefaultsecuritycontext-ref-gwp-on-gw",
		},
	}

	tester := DeployerTester{
		ControllerName:    wellknown.DefaultGatewayControllerName,
		AgwControllerName: wellknown.DefaultAgwControllerName,
		ClassName:         wellknown.DefaultGatewayClassName,
		WaypointClassName: wellknown.DefaultWaypointClassName,
		AgwClassName:      wellknown.DefaultAgwClassName,
	}

	dir := fsutils.MustGetThisDir()
	scheme := schemes.GatewayScheme()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tester.RunHelmChartTest(t, tt, scheme, dir, nil)
		})
	}
}

func TestRenderHelmChartWithTLS(t *testing.T) {
	// Create temporary CA certificate file for testing
	caCertContent := `-----BEGIN CERTIFICATE-----
MIICljCCAX4CCQCKSGhvPtMNGzANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJV
UzAeFw0yNDA3MDEwMDAwMDBaFw0yNTA3MDEwMDAwMDBaMA0xCzAJBgNVBAYTAlVT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234567890ABCDEFGHIj
klmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab
cdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ123456
7890abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ
1234567890abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTU
VWXYZ1234567890abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNO
PQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHI
JKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz1234567890ABC
DEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz123456
wIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBtestcertdata
-----END CERTIFICATE-----`

	// Create temporary directory and CA certificate file
	tmpDir := t.TempDir()
	caCertPath := tmpDir + "/ca.crt"
	err := os.WriteFile(caCertPath, []byte(caCertContent), 0o600)
	require.NoError(t, err)

	tests := []HelmTestCase{
		{
			Name:      "basic gateway with TLS enabled",
			InputFile: "base-gateway-tls",
		},
		{
			Name:      "agentgateway with TLS enabled",
			InputFile: "agentgateway-tls",
		},
	}

	tester := DeployerTester{
		ControllerName:    wellknown.DefaultGatewayControllerName,
		AgwControllerName: wellknown.DefaultAgwControllerName,
		ClassName:         wellknown.DefaultGatewayClassName,
		WaypointClassName: wellknown.DefaultWaypointClassName,
		AgwClassName:      wellknown.DefaultAgwClassName,
	}

	// ExtraGatewayParameters function that enables TLS. this is needed as TLS
	// is injected by the control plane and not via the GWP API.
	tlsExtraParams := func(cli client.Client, inputs *pkgdeployer.Inputs) pkgdeployer.HelmValuesGenerator {
		inputs.ControlPlane.XdsTLS = true
		inputs.ControlPlane.XdsTlsCaPath = caCertPath
		return nil
	}

	dir := fsutils.MustGetThisDir()
	scheme := schemes.GatewayScheme()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tester.RunHelmChartTest(t, tt, scheme, dir, tlsExtraParams)
		})
	}
}
