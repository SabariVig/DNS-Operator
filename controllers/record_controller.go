package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	domainv1 "github.com/SabariVig/DNS-Operator/api/v1"
	v1 "github.com/SabariVig/DNS-Operator/api/v1"
	"github.com/SabariVig/DNS-Operator/pkg/namecheap"
	"github.com/SabariVig/DNS-Operator/pkg/providers"
	corev1 "k8s.io/api/core/v1"
)

// RecordReconciler reconciles a Record object
type RecordReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Secret corev1.Secret
}

//+kubebuilder:rbac:groups=domain.lxz.io,resources=records,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=domain.lxz.io,resources=records/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=domain.lxz.io,resources=records/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RecordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Here")
	var provider providers.Providers
	var record v1.Record
	err := r.Get(ctx, req.NamespacedName, &record)
	if err != nil {
		log.Error(err, "Unable to get record from reconciler")
		return ctrl.Result{}, err
	}

	log.Info("Provider", "Provider", record.Spec.Provider)

	switch record.Spec.Provider {
	case v1.Namecheap:
		provider = namecheap.NewNamecheapProvider(r.Secret)
	}

	finalizerName := "records.domain.lxz.io/finalizers"

	if record.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Finalizers", "Deletion Issued")
		if !controllerutil.ContainsFinalizer(&record, finalizerName) {
			log.Info("Finalizers: Adding finalizer to", "hostname", &record.Spec.Hostname)

			controllerutil.AddFinalizer(&record, finalizerName)
			err := r.Update(ctx, &record)
			if err != nil {
				log.Error(err, "Unable to add Finalizers")
				return ctrl.Result{}, err
			}
			log.Info("Finalizers: Added finalizer to", "hostname", &record.Spec.Hostname)
		}
	} else {
		if controllerutil.ContainsFinalizer(&record, finalizerName) {

			log.Info("Finalizers", "Deletion Started")
			err := provider.DeleteRecord(ctx, &record)
			if err != nil {
				log.Error(err, "Unable to Delete Record")
				return ctrl.Result{}, err
			}

			log.Info("Finalizers: Removing Finalizers from", "record", record)
			controllerutil.RemoveFinalizer(&record, finalizerName)
			err = r.Update(ctx, &record)
			if err != nil {
				log.Error(err, "Unable to remove Finalizers")
				return ctrl.Result{}, err
			}
			log.Info("Finalizers: removed Finalizers from", "record", record)
		}
		return ctrl.Result{}, err
	}

	err = provider.AddRecord(ctx, &record)
	if err != nil {
		log.Error(err, "Unable to add Record")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RecordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&domainv1.Record{}).
		Complete(r)
}
