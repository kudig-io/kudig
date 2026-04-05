// Package controllers contains controller implementations for kudig CRDs
package controllers

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kudigv1 "github.com/kudig/kudig/operator/api/v1"
)

// ClusterDiagnosticReconciler reconciles a ClusterDiagnostic object
type ClusterDiagnosticReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kudig.io,resources=clusterdiagnostics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kudig.io,resources=clusterdiagnostics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kudig.io,resources=clusterdiagnostics/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main reconciliation loop
func (r *ClusterDiagnosticReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ClusterDiagnostic instance
	diagnostic := &kudigv1.ClusterDiagnostic{}
	if err := r.Get(ctx, req.NamespacedName, diagnostic); err != nil {
		if errors.IsNotFound(err) {
			log.Info("ClusterDiagnostic resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ClusterDiagnostic")
		return ctrl.Result{}, err
	}

	// Check if already completed
	if diagnostic.Status.Phase == "Completed" || diagnostic.Status.Phase == "Failed" {
		return ctrl.Result{}, nil
	}

	// Check if there's an active job
	jobName := fmt.Sprintf("kudig-cluster-%s", diagnostic.Name)
	found := &batchv1.Job{}
	err := r.Get(ctx, types.NamespacedName{Name: jobName, Namespace: "kudig-system"}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create a new job
		job := r.jobForClusterDiagnostic(diagnostic, jobName)
		log.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		if err := r.Create(ctx, job); err != nil {
			log.Error(err, "Failed to create new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
			return ctrl.Result{}, err
		}

		// Update status to Running
		diagnostic.Status.Phase = "Running"
		diagnostic.Status.StartTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update ClusterDiagnostic status")
			return ctrl.Result{}, err
		}

		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Job")
		return ctrl.Result{}, err
	}

	// Job exists, check its status
	if found.Status.Succeeded > 0 {
		// Job completed successfully
		diagnostic.Status.Phase = "Completed"
		now := metav1.Now()
		diagnostic.Status.CompletionTime = &now

		// TODO: Parse job output to populate summary
		diagnostic.Status.Summary = kudigv1.DiagnosticSummary{
			Total:        0,
			Critical:     0,
			Warning:      0,
			Info:         0,
			AnalyzersRun: len(diagnostic.Spec.Analyzers),
		}

		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update ClusterDiagnostic status")
			return ctrl.Result{}, err
		}
	} else if found.Status.Failed > 0 {
		// Job failed
		diagnostic.Status.Phase = "Failed"
		now := metav1.Now()
		diagnostic.Status.CompletionTime = &now

		condition := kudigv1.DiagnosticCondition{
			Type:               "Failed",
			Status:             "True",
			LastTransitionTime: now,
			Reason:             "JobFailed",
			Message:            "The diagnostic job failed",
		}
		diagnostic.Status.Conditions = append(diagnostic.Status.Conditions, condition)

		if err := r.Status().Update(ctx, diagnostic); err != nil {
			log.Error(err, "Failed to update ClusterDiagnostic status")
			return ctrl.Result{}, err
		}
	} else {
		// Job still running
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

// jobForClusterDiagnostic creates a Job for a ClusterDiagnostic resource
func (r *ClusterDiagnosticReconciler) jobForClusterDiagnostic(diagnostic *kudigv1.ClusterDiagnostic, jobName string) *batchv1.Job {
	labels := map[string]string{
		"app":                    "kudig",
		"kudig.io/diagnostic":    diagnostic.Name,
		"kudig.io/type":          "cluster",
		"kudig.io/created-by":    "operator",
	}

	// Build command arguments
	args := []string{"online", "--format", diagnostic.Spec.OutputFormat}
	if len(diagnostic.Spec.Analyzers) > 0 {
		for _, analyzer := range diagnostic.Spec.Analyzers {
			args = append(args, "--analyzer", analyzer)
		}
	}
	if len(diagnostic.Spec.ExcludeAnalyzers) > 0 {
		for _, analyzer := range diagnostic.Spec.ExcludeAnalyzers {
			args = append(args, "--exclude-analyzer", analyzer)
		}
	}

	backoffLimit := int32(2)
	ttlSecondsAfterFinished := int32(3600) // 1 hour

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
								Name: "KUBERNETES_NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
						},
					}},
					ServiceAccountName: "kudig-operator",
				},
			},
		},
	}

	// Set owner reference
	ctrl.SetControllerReference(diagnostic, job, r.Scheme)

	return job
}

// SetupWithManager sets up the controller with the Manager
func (r *ClusterDiagnosticReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kudigv1.ClusterDiagnostic{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
