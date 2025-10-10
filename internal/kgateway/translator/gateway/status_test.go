package gateway_test

import (
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
	translatortest "github.com/kgateway-dev/kgateway/v2/test/translator"
)

func TestStatuses(t *testing.T) {
	testFn := func(t *testing.T, inputFile string, wantPolicyErrors map[reporter.PolicyKey]*gwv1.PolicyStatus) {
		dir := fsutils.MustGetThisDir()
		translatortest.TestTranslation(
			t,
			t.Context(),
			[]string{
				filepath.Join(dir, "testutils/inputs/status", inputFile),
			},
			filepath.Join(dir, "testutils/outputs/status", inputFile),
			types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		)
	}

	t.Run("Basic", func(t *testing.T) {
		testFn(t, "basic.yaml", map[reporter.PolicyKey]*gwv1.PolicyStatus{
			{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "default", Name: "extensionref-policy"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 1,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 1,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonMerged),
								Message:            reporter.PolicyMergedMsg,
							},
						},
					},
				},
			},
			{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "default", Name: "policy-with-section-name"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 2,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 2,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonMerged),
								Message:            reporter.PolicyMergedMsg,
							},
						},
					},
				},
			},
			{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "default", Name: "policy-without-section-name"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 3,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 3,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonMerged),
								Message:            reporter.PolicyMergedMsg,
							},
						},
					},
				},
			},
			{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "default", Name: "fully-ignored"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 4,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 4,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionFalse,
								Reason:             string(v1alpha1.PolicyReasonOverridden),
								Message:            reporter.PolicyOverriddenMsg,
							},
						},
					},
				},
			},
			{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "default", Name: "policy-no-merge"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 1,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 1,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonAttached),
								Message:            reporter.PolicyAttachedMsg,
							},
						},
					},
				},
			},
			{Group: "gateway.kgateway.dev", Kind: "HTTPListenerPolicy", Namespace: "default", Name: "policy-1"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 1,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 1,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonMerged),
								Message:            reporter.PolicyMergedMsg,
							},
						},
					},
				},
			},
			{Group: "gateway.kgateway.dev", Kind: "HTTPListenerPolicy", Namespace: "default", Name: "policy-2"}: {
				Ancestors: []gwv1.PolicyAncestorStatus{
					{
						AncestorRef: gwv1.ParentReference{
							Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
							Kind:      ptr.To(gwv1.Kind("Gateway")),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Name:      gwv1.ObjectName("example-gateway"),
						},
						ControllerName: wellknown.DefaultGatewayControllerName,
						Conditions: []metav1.Condition{
							{
								ObservedGeneration: 2,
								Type:               string(v1alpha1.PolicyConditionAccepted),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonValid),
								Message:            reporter.PolicyAcceptedMsg,
							},
							{
								ObservedGeneration: 2,
								Type:               string(v1alpha1.PolicyConditionAttached),
								Status:             metav1.ConditionTrue,
								Reason:             string(v1alpha1.PolicyReasonMerged),
								Message:            reporter.PolicyMergedMsg,
							},
						},
					},
				},
			},
		})
	})
}
