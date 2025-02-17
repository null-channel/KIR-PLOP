/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	treev1alpha1 "github.com/null-channel/job-hunting/tree-operator/api/v1alpha1"
)

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = treev1alpha1.GroupVersion.String()
)

// TreeReconciler reconciles a Tree object
type TreeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tree.nullcloud.io,resources=trees,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tree.nullcloud.io,resources=trees/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tree.nullcloud.io,resources=trees/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tree object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *TreeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	// WOOT! We have a Tree object. Let's do something with it.

	// Fetch the Tree instance
	tree := &treev1alpha1.Tree{}
	if err := r.Get(ctx, req.NamespacedName, tree); err != nil {
		log.Log.Error(err, "unable to fetch Tree")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// your logic here
	// Get all the pods that are owned by this operator

	pods := &corev1.PodList{}
	if err := r.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Log.Error(err, "unable to list child Pods")
		return ctrl.Result{}, err
	}

	// Check if depth is equal to the number of pods
	// If not, create the pods
	if len(pods.Items) != tree.Spec.Count {
		// Delete all old pods
		// TODO: This is not the best way to do this
		// We should be checking if the pods are still needed
		// and only delete the ones that are not needed
		// or add the ones needed.
		for _, pod := range pods.Items {
			if err := r.Delete(ctx, &pod); err != nil {
				log.Log.Error(err, "unable to delete old pods")
				return ctrl.Result{}, err
			}
		}
		var root *Node

		for i := 1; i <= tree.Spec.Count; i++ {
			root = Insert(root, i)
		}

		// Create new pods
		iterator := NewLevelOrderIterator(root)
		for node := iterator.Next(); node != nil; node = iterator.Next() {
			labels := map[string]string{}
			if node.Right != nil {
				labels["r-child"] = fmt.Sprintf("%d", node.Right.Key)
			}
			if node.Left != nil {
				labels["l-child"] = fmt.Sprintf("%d", node.Left.Key)
			}
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%d", tree.Name, node.Key),
					Namespace: tree.Namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: tree.APIVersion,
							Kind:       tree.Kind,
							Name:       tree.Name,
							UID:        tree.UID,
						},
					},
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "tree",
							Image: "nginx",
						},
					},
				},
			}
			if err := ctrl.SetControllerReference(tree, pod, r.Scheme); err != nil {
				log.Log.Error(err, "unable to set owner reference on pod")
				return ctrl.Result{}, err
			}
			if err := r.Create(ctx, pod); err != nil {
				log.Log.Error(err, "unable to create pod")
				return ctrl.Result{}, err
			}
		}

	}
	// Check to see if we have all the pods we need for the tree
	// and that they all have a referance to their children

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TreeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Pod{}, jobOwnerKey, func(rawObj client.Object) []string {
		// grab the tree object, extract the owner...
		pod := rawObj.(*corev1.Pod)
		owner := metav1.GetControllerOf(pod)
		if owner == nil {
			return nil
		}

		// ...make sure it's a Tree...
		if owner.APIVersion != apiGVStr || owner.Kind != "Tree" {
			return nil
		}

		fmt.Println(owner.Name)

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&treev1alpha1.Tree{}).
		Named("tree").
		Complete(r)
}
