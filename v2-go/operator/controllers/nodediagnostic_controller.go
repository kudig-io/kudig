package controllers

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kudigv1 "github.com/kudig/kudig/operator/api/v1"
)

// NodeDiagnosticReconciler reconciles a NodeDiagnostic object
type NodeDiagnosticReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kudig.io,resources=nodediagnostics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kudig.io,resources=nodediagnostics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kudig.io,resources=nodediagnostics/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch

// Reconcile is the main reconciliation loop
func (r *NodeDiagnosticReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the NodeDiagnostic instance
	diagnostic := &kudigv1.NodeDiagnostic{}
	if err := r.Get(ctx, req.NamespacedName, diagnostic); err != nil {
		if errors.IsNotFound(err) {
			log.Info("NodeDiagnostic resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get NodeDiagnostic")
		return ctrl.Result{}, err
	}

	// Check if already completed
	if diagnostic.Status.Phase == "Completed" || diagnostic.Status.Phase == "Failed" {
		return ctrl.Result{}, nil
	}

	// Get target nodes
	targetNodes, err := r.getTargetNodes(ctx, diagnostic)
	if err != nil {
		log.Error(err, "Failed to get target nodes")
		return ctrl.Result{}, err
	}

	if len(targetNodes) == 0 {
		log.Info("No target nodes found for diagnostic")
		diagnostic.Status.Phase = "Failed"
		now := metav1.Now()
		diagnostic.Status.CompletionTime = &now
		condition := kudigv1.DiagnosticCondition{
			Type:               "Failed",
			Status:             "True",
			LastTransitionTime: now,
			Reason:             "NoNodes",
			Message:            "No target nodes found",
		}
		diagnostic.Status.Conditions = append(diagnostic.Status.Conditions, condition)
		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update NodeDiagnostic status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Initialize status if not set
	if diagnostic.Status.Phase == "" {
		diagnostic.Status.Phase = "Running"
		diagnostic.Status.StartTime = &metav1.Time{Time: time.Now()}
		diagnostic.Status.NodeResults = make([]kudigv1.NodeResult, 0, len(targetNodes))

		for _, node := range targetNodes {
			diagnostic.Status.NodeResults = append(diagnostic.Status.NodeResults, kudigv1.NodeResult{
				NodeName: node,
				Phase:    "Pending",
			})
		}

		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update NodeDiagnostic status")
			return ctrl.Result{}, err
		}
	}

	// Create DaemonSet to run diagnostics on each node
	dsName := fmt.Sprintf("kudig-node-%s", diagnostic.Name)
	found := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: dsName, Namespace: "kudig-system"}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create DaemonSet-like Job for node diagnostics
		// For simplicity, we create a single job that runs on all nodes using host networking
		job := r.jobForNodeDiagnostic(diagnostic, dsName, targetNodes)
		log.Info("Creating a new Job for node diagnostics", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		if err := r.Create(ctx, job); err != nil {
			log.Error(err, "Failed to create new Job")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Job")
		return ctrl.Result{}, err
	}

	// Check job status
	if found.Status.Succeeded > 0 {
		diagnostic.Status.Phase = "Completed"
		now := metav1.Now()
		diagnostic.Status.CompletionTime = &now

		// Update all node results
		for i := range diagnostic.Status.NodeResults {
			diagnostic.Status.NodeResults[i].Phase = "Completed"
		}

		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update NodeDiagnostic status")
			return ctrl.Result{}, err
		}
	} else if found.Status.Failed > 0 {
		diagnostic.Status.Phase = "Failed"
		now := metav1.Now()
		diagnostic.Status.CompletionTime = &now

		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update NodeDiagnostic status")
			return ctrl.Result{}, err
		}
	} else {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

// getTargetNodes returns the list of nodes to diagnose
func (r *NodeDiagnosticReconciler) getTargetNodes(ctx context.Context, diagnostic *kudigv1.NodeDiagnostic) ([]string, error) {
	// If specific node names are provided, use them
	if len(diagnostic.Spec.NodeNames) > 0 {
		return diagnostic.Spec.NodeNames, nil
	}

	// Otherwise, use node selector
	nodeList := &corev1.NodeList{}
	selector := labels.Everything()
	if len(diagnostic.Spec.NodeSelector) > 0 {
		var err error
		selector, err = metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
			MatchLabels: diagnostic.Spec.NodeSelector,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := r.List(ctx, nodeList, &client.ListOptions{LabelSelector: selector}); err != nil {
		return nil, err
	}

	nodes := make([]string, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		nodes = append(nodes, node.Name)
	}

	return nodes, nil
}

// jobForNodeDiagnostic creates a Job for a NodeDiagnostic resource
func (r *NodeDiagnosticReconciler) jobForNodeDiagnostic(diagnostic *kudigv1.NodeDiagnostic, jobName string, nodes []string) *batchv1.Job {
	labels := map[string]string{
		"app":                 "kudig",
		"kudig.io/diagnostic": diagnostic.Name,
		"kudig.io/type":       "node",
		"kudig.io/created-by": "operator",
	}

	// Build command arguments
	args := []string{"online", "--all-nodes", "--format", diagnostic.Spec.OutputFormat}
	if len(diagnostic.Spec.Analyzers) > 0 {
		for _, analyzer := range diagnostic.Spec.Analyzers {
			args = append(args, "--analyzer", analyzer)
		}
	}

	backoffLimit := int32(2)
	ttlSecondsAfterFinished := int32(3600)

	// Add node names as env var for reference
	nodeNamesEnv := ""
	for i, node := range nodes {
		if i > 0 {
			nodeNamesEnv += ","
		}
		nodeNamesEnv += node
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: "kudig-system",
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{{
						Name:    "kudig",
						Image:   "kudig/kudig:latest",
						Command: []string{"/usr/local/bin/kudig"},
						Args:    args,
						Env: []corev1.EnvVar{
							{
								Name:  "TARGET_NODES",
								Value: nodeNamesEnv,
							},
						},
					}},
					ServiceAccountName: "kudig-operator",
				},
			},
		},
	}

	ctrl.SetControllerReference(diagnostic, job, r.Scheme)

	return job
}

// SetupWithManager sets up the controller with the Manager
func (r *NodeDiagnosticReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kudigv1.NodeDiagnostic{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
